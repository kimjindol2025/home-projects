#!/usr/bin/env python3
"""
Phase 1: 통합 빌드 스캔
130개 프로젝트의 메타데이터 수집 및 빌드 규칙 분석

출력:
- build_inventory.csv: 프로젝트별 빌드 메타데이터
- dependency_graph.csv: 프로젝트 간 의존성
- build_summary.md: 분석 리포트
"""

import os
import json
import csv
from pathlib import Path
from collections import defaultdict
import re
from datetime import datetime

class BuildScanner:
    def __init__(self, root_path):
        self.root_path = Path(root_path)
        self.projects = []
        self.build_files = {
            'Cargo.toml': 'rust',
            'package.json': 'node',
            'Makefile': 'make',
            'CMakeLists.txt': 'cmake',
            'pyproject.toml': 'python',
            'setup.py': 'python',
            'build.gradle': 'gradle',
            'pom.xml': 'maven',
        }
        self.freelang_extensions = {'.fl', '.fv'}

    def scan(self):
        """모든 프로젝트 스캔"""
        print("🔍 프로젝트 스캔 시작...")

        # .projects/core, archived, experiments, modules
        for category_dir in ['core', 'archived', 'experiments', 'modules']:
            category_path = self.root_path / '.projects' / category_dir
            if category_path.exists():
                self._scan_category(category_path, category_dir)

        # projects/ 디렉토리도 스캔
        projects_path = self.root_path / 'projects'
        if projects_path.exists():
            self._scan_category(projects_path, 'top_level')

        print(f"✅ 총 {len(self.projects)}개 프로젝트 발견")
        return self.projects

    def _scan_category(self, category_path, category_name):
        """카테고리별 프로젝트 스캔"""
        for project_dir in category_path.iterdir():
            if not project_dir.is_dir() or project_dir.name.startswith('.'):
                continue

            project_info = self._analyze_project(project_dir, category_name)
            if project_info:
                self.projects.append(project_info)

    def _analyze_project(self, project_path, category):
        """프로젝트 분석"""
        project_name = project_path.name

        # 빌드 파일 감지
        build_system = self._detect_build_system(project_path)

        # FreeLang 파일 검색
        fl_files = self._find_freelang_files(project_path)

        # 의존성 검색
        dependencies = self._extract_dependencies(project_path, build_system)

        # 소스 통계
        stats = self._count_files(project_path)

        return {
            'name': project_name,
            'category': category,
            'path': str(project_path.relative_to(self.root_path)),
            'build_system': build_system,
            'has_freelang': len(fl_files) > 0,
            'fl_file_count': len(fl_files),
            'fl_files': ','.join(fl_files[:5]),  # 처음 5개만
            'dependencies': ','.join(dependencies[:10]),  # 처음 10개만
            'src_files': stats['src'],
            'test_files': stats['test'],
            'total_lines': stats['lines'],
            'git_commits': self._count_commits(project_path),
        }

    def _detect_build_system(self, project_path):
        """빌드 시스템 감지"""
        for build_file, build_type in self.build_files.items():
            if (project_path / build_file).exists():
                return build_type
        return 'unknown'

    def _find_freelang_files(self, project_path):
        """FreeLang 파일 검색"""
        fl_files = []
        for root, dirs, files in os.walk(project_path):
            # 깊이 제한 (성능)
            depth = root[len(str(project_path)):].count(os.sep)
            if depth > 3:
                dirs.clear()
                continue

            for file in files:
                if Path(file).suffix in self.freelang_extensions:
                    fl_files.append(file)
        return fl_files[:20]  # 최대 20개

    def _extract_dependencies(self, project_path, build_system):
        """의존성 추출"""
        deps = []

        if build_system == 'rust':
            cargo_path = project_path / 'Cargo.toml'
            if cargo_path.exists():
                with open(cargo_path) as f:
                    content = f.read()
                    # [dependencies] 섹션 파싱
                    in_deps = False
                    for line in content.split('\n'):
                        if '[dependencies]' in line:
                            in_deps = True
                        elif in_deps and line.startswith('['):
                            break
                        elif in_deps and '=' in line:
                            dep = line.split('=')[0].strip()
                            if dep:
                                deps.append(dep)

        elif build_system == 'node':
            pkg_path = project_path / 'package.json'
            if pkg_path.exists():
                try:
                    with open(pkg_path) as f:
                        pkg = json.load(f)
                        deps = list(pkg.get('dependencies', {}).keys())[:20]
                except:
                    pass

        return deps

    def _count_files(self, project_path):
        """파일 통계"""
        stats = {'src': 0, 'test': 0, 'lines': 0}

        for root, dirs, files in os.walk(project_path):
            # 깊이 제한
            depth = root[len(str(project_path)):].count(os.sep)
            if depth > 4:
                dirs.clear()
                continue

            for file in files:
                file_path = Path(root) / file

                # 소스 파일인지 확인
                if any(file.endswith(ext) for ext in ['.fl', '.fv', '.rs', '.py', '.js', '.c', '.h']):
                    if 'test' in root:
                        stats['test'] += 1
                    else:
                        stats['src'] += 1

                    # 줄 수 계산
                    try:
                        with open(file_path) as f:
                            stats['lines'] += len(f.readlines())
                    except:
                        pass

        return stats

    def _count_commits(self, project_path):
        """Git 커밋 수 계산"""
        try:
            git_dir = project_path / '.git'
            if git_dir.exists():
                result = os.popen(f'cd {project_path} && git rev-list --count HEAD 2>/dev/null').read().strip()
                return int(result) if result else 0
        except:
            pass
        return 0


class DependencyAnalyzer:
    """프로젝트 간 의존성 분석"""

    def __init__(self, projects):
        self.projects = projects
        self.project_map = {p['name']: p for p in projects}

    def analyze(self):
        """의존성 그래프 생성"""
        dependencies = []

        for project in self.projects:
            deps = project.get('dependencies', '').split(',')
            for dep in deps:
                dep = dep.strip()
                if dep and dep in self.project_map:
                    dependencies.append({
                        'from': project['name'],
                        'to': dep,
                        'type': 'internal',
                    })

        return dependencies


def export_csv(projects, dependencies):
    """CSV로 내보내기"""

    # build_inventory.csv
    print("📊 build_inventory.csv 생성...")

    # 필드 제거 (dependencies 제거)
    clean_projects = []
    for p in projects:
        clean_p = {k: v for k, v in p.items() if k != 'dependencies'}
        clean_projects.append(clean_p)

    with open('build_inventory.csv', 'w', newline='') as f:
        writer = csv.DictWriter(f, fieldnames=[
            'name', 'category', 'path', 'build_system',
            'has_freelang', 'fl_file_count', 'fl_files',
            'src_files', 'test_files', 'total_lines', 'git_commits'
        ])
        writer.writeheader()
        writer.writerows(clean_projects)

    # dependency_graph.csv
    print("📊 dependency_graph.csv 생성...")
    with open('dependency_graph.csv', 'w', newline='') as f:
        writer = csv.DictWriter(f, fieldnames=['from', 'to', 'type'])
        writer.writeheader()
        writer.writerows(dependencies)


def generate_report(projects, dependencies):
    """분석 리포트 생성"""

    report = []
    report.append("# Phase 1: 통합 빌드 스캔 리포트\n")
    report.append(f"생성일: {datetime.now().isoformat()}\n\n")

    # 개요
    report.append("## 📊 개요\n")
    report.append(f"- 총 프로젝트: {len(projects)}\n")
    report.append(f"- FreeLang 프로젝트: {sum(1 for p in projects if p['has_freelang'])}\n")
    report.append(f"- 내부 의존성: {len(dependencies)}\n")

    # 빌드 시스템 분포
    report.append("\n## 🔨 빌드 시스템 분포\n")
    build_systems = defaultdict(int)
    for p in projects:
        build_systems[p['build_system']] += 1

    for system, count in sorted(build_systems.items(), key=lambda x: -x[1]):
        report.append(f"- {system}: {count}\n")

    # 카테고리별 분포
    report.append("\n## 📂 카테고리별 분포\n")
    categories = defaultdict(int)
    for p in projects:
        categories[p['category']] += 1

    for cat, count in sorted(categories.items()):
        report.append(f"- {cat}: {count}\n")

    # FreeLang 프로젝트
    report.append("\n## 🎯 FreeLang 프로젝트\n")
    fl_projects = [p for p in projects if p['has_freelang']]
    report.append(f"총 {len(fl_projects)}개\n\n")

    for p in sorted(fl_projects, key=lambda x: -x['fl_file_count'])[:20]:
        report.append(f"- **{p['name']}**: {p['fl_file_count']}개 파일, {p['total_lines']}줄\n")

    # 상위 프로젝트 (코드 규모)
    report.append("\n## 📈 상위 프로젝트 (코드 규모)\n")
    for p in sorted(projects, key=lambda x: -x['total_lines'])[:15]:
        report.append(f"- **{p['name']}**: {p['total_lines']}줄, {p['git_commits']} commits\n")

    # 의존성 네트워크
    report.append("\n## 🔗 의존성 네트워크\n")
    in_degree = defaultdict(int)
    out_degree = defaultdict(int)

    for dep in dependencies:
        out_degree[dep['from']] += 1
        in_degree[dep['to']] += 1

    report.append("\n### 가장 많이 의존되는 프로젝트\n")
    for proj, count in sorted(in_degree.items(), key=lambda x: -x[1])[:10]:
        report.append(f"- {proj}: {count}개 프로젝트에서 의존\n")

    report.append("\n### 가장 많이 의존하는 프로젝트\n")
    for proj, count in sorted(out_degree.items(), key=lambda x: -x[1])[:10]:
        report.append(f"- {proj}: {count}개 프로젝트에 의존\n")

    return ''.join(report)


def main():
    root = Path('/data/data/com.termux/files/home')

    # 스캔
    scanner = BuildScanner(root)
    projects = scanner.scan()

    # 의존성 분석
    analyzer = DependencyAnalyzer(projects)
    dependencies = analyzer.analyze()

    # 내보내기
    export_csv(projects, dependencies)

    # 리포트 생성
    report = generate_report(projects, dependencies)
    with open('build_summary.md', 'w') as f:
        f.write(report)

    print("\n✅ Phase 1 완료!")
    print("📁 생성된 파일:")
    print("  - build_inventory.csv")
    print("  - dependency_graph.csv")
    print("  - build_summary.md")


if __name__ == '__main__':
    main()
