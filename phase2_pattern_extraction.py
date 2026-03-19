#!/usr/bin/env python3
"""
Phase 2: 공통 모듈 추출
130개+ 프로젝트에서 반복되는 로직 식별 및 Standard Core 라이브러리 설계

전략:
1. 모든 FreeLang 파일 분석
2. 함수/모듈 패턴 추출
3. 반복도(frequency) 계산
4. Standard Core 모듈 구성도 설계
"""

import os
import re
from pathlib import Path
from collections import defaultdict, Counter
import json
from datetime import datetime

class PatternExtractor:
    def __init__(self, root_path):
        self.root_path = Path(root_path)
        self.patterns = defaultdict(int)
        self.modules = defaultdict(list)
        self.function_defs = defaultdict(int)
        self.imports = defaultdict(int)
        self.all_files = []

    def extract_all(self):
        """모든 FreeLang 파일 분석"""
        print("🔍 Phase 2: 패턴 추출 시작...\n")

        # .projects 디렉토리 스캔
        for category_dir in ['core', 'archived', 'experiments', 'modules']:
            category_path = self.root_path / '.projects' / category_dir
            if category_path.exists():
                self._scan_fl_files(category_path)

        # projects/ 디렉토리 스캔
        projects_path = self.root_path / 'projects'
        if projects_path.exists():
            self._scan_fl_files(projects_path)

        print(f"✅ 총 {len(self.all_files)}개 FreeLang 파일 분석 완료\n")
        return self.all_files

    def _scan_fl_files(self, directory):
        """FreeLang 파일 스캔 및 분석"""
        for root, dirs, files in os.walk(directory):
            depth = root[len(str(directory)):].count(os.sep)
            if depth > 4:
                dirs.clear()
                continue

            for file in files:
                if file.endswith(('.fl', '.fv')):
                    file_path = Path(root) / file
                    try:
                        self._analyze_file(file_path)
                    except Exception as e:
                        pass

    def _analyze_file(self, file_path):
        """단일 파일 분석"""
        with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()

        file_info = {
            'path': str(file_path),
            'name': file_path.name,
            'lines': len(content.split('\n')),
            'functions': self._extract_functions(content),
            'modules': self._extract_modules(content),
            'imports': self._extract_imports(content),
            'patterns': self._extract_patterns(content),
        }

        self.all_files.append(file_info)

        # 통계 업데이트
        for func in file_info['functions']:
            self.function_defs[func] += 1
        for mod in file_info['modules']:
            self.modules[mod].append(file_path.name)
        for imp in file_info['imports']:
            self.imports[imp] += 1
        for pat in file_info['patterns']:
            self.patterns[pat] += 1

    def _extract_functions(self, content):
        """함수 정의 추출"""
        functions = []
        # fn name(...) 패턴
        pattern = r'fn\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\('
        for match in re.finditer(pattern, content):
            functions.append(match.group(1))
        return functions

    def _extract_modules(self, content):
        """모듈 정의 추출"""
        modules = []
        # mod name { ... } 패턴
        pattern = r'mod\s+([a-zA-Z_][a-zA-Z0-9_]*)'
        for match in re.finditer(pattern, content):
            modules.append(match.group(1))
        return modules

    def _extract_imports(self, content):
        """import 문 추출"""
        imports = []
        # use ... 또는 import ... 패턴
        pattern = r'(?:use|import)\s+([a-zA-Z_][a-zA-Z0-9_:]*)'
        for match in re.finditer(pattern, content):
            imports.append(match.group(1))
        return imports

    def _extract_patterns(self, content):
        """디자인 패턴 식별"""
        patterns = []

        # 자주 사용되는 패턴들
        if re.search(r'impl\s+\w+\s*\{', content):
            patterns.append('impl_block')
        if re.search(r'struct\s+\w+\s*\{', content):
            patterns.append('struct')
        if re.search(r'enum\s+\w+\s*\{', content):
            patterns.append('enum')
        if re.search(r'trait\s+\w+\s*\{', content):
            patterns.append('trait')
        if re.search(r'Result<', content):
            patterns.append('result_pattern')
        if re.search(r'Option<', content):
            patterns.append('option_pattern')
        if re.search(r'match\s+\w+\s*\{', content):
            patterns.append('pattern_matching')
        if re.search(r'loop\s*\{', content):
            patterns.append('loop')
        if re.search(r'for\s+\w+\s+in\s+', content):
            patterns.append('for_loop')
        if re.search(r'while\s+', content):
            patterns.append('while_loop')
        if re.search(r'async\s+fn', content):
            patterns.append('async_function')
        if re.search(r'await\b', content):
            patterns.append('await')
        if re.search(r'#\[derive\(', content):
            patterns.append('derive_macro')
        if re.search(r'#\[test\]', content):
            patterns.append('unit_test')
        if re.search(r'unsafe\s*\{', content):
            patterns.append('unsafe_block')
        if re.search(r'pub\s+fn', content):
            patterns.append('public_function')
        if re.search(r'pub\s+struct', content):
            patterns.append('public_struct')
        if re.search(r'pub\s+trait', content):
            patterns.append('public_trait')

        return patterns


class StandardCoreDesigner:
    """Standard Core 라이브러리 설계"""

    def __init__(self, extractor):
        self.extractor = extractor

    def design(self):
        """Standard Core 구성도 설계"""

        core_modules = {
            'collections': {
                'functions': ['Vec', 'HashMap', 'HashSet', 'LinkedList', 'BTreeMap'],
                'frequency': 0,
                'description': '기본 자료구조',
            },
            'string': {
                'functions': ['String', 'str', 'trim', 'split', 'join', 'format'],
                'frequency': 0,
                'description': '문자열 처리',
            },
            'io': {
                'functions': ['File', 'read', 'write', 'println', 'BufReader', 'BufWriter'],
                'frequency': 0,
                'description': '입출력',
            },
            'math': {
                'functions': ['sqrt', 'sin', 'cos', 'pow', 'abs', 'max', 'min'],
                'frequency': 0,
                'description': '수학 연산',
            },
            'result': {
                'functions': ['Result', 'Ok', 'Err', 'map', 'unwrap', 'expect'],
                'frequency': 0,
                'description': '에러 처리',
            },
            'option': {
                'functions': ['Option', 'Some', 'None', 'map', 'unwrap', 'unwrap_or'],
                'frequency': 0,
                'description': '선택적 값',
            },
            'iter': {
                'functions': ['Iterator', 'map', 'filter', 'fold', 'zip', 'enumerate'],
                'frequency': 0,
                'description': '반복자',
            },
            'async': {
                'functions': ['async', 'await', 'Future', 'Executor', 'spawn'],
                'frequency': 0,
                'description': '비동기',
            },
            'concurrency': {
                'functions': ['Mutex', 'RwLock', 'Channel', 'sync', 'thread'],
                'frequency': 0,
                'description': '동시성',
            },
            'serialization': {
                'functions': ['serialize', 'deserialize', 'json', 'serde'],
                'frequency': 0,
                'description': '직렬화',
            },
            'crypto': {
                'functions': ['sha256', 'md5', 'encrypt', 'decrypt', 'hash'],
                'frequency': 0,
                'description': '암호화',
            },
            'network': {
                'functions': ['TcpStream', 'TcpListener', 'http', 'request'],
                'frequency': 0,
                'description': '네트워크',
            },
        }

        # 실제 사용 빈도 계산
        for module_name, module_info in core_modules.items():
            for func in module_info['functions']:
                if func in self.extractor.function_defs:
                    module_info['frequency'] += self.extractor.function_defs[func]

        return core_modules

    def estimate_loc(self):
        """Standard Core 예상 코드량"""
        estimates = {
            'collections': 400,
            'string': 300,
            'io': 250,
            'math': 150,
            'result': 120,
            'option': 120,
            'iter': 300,
            'async': 400,
            'concurrency': 350,
            'serialization': 400,
            'crypto': 450,
            'network': 500,
        }
        return estimates


def generate_report(extractor, core_modules, estimates):
    """상세 분석 리포트 생성"""

    report = []
    report.append("# Phase 2: 공통 모듈 추출 & Standard Core 설계\n\n")
    report.append(f"생성일: {datetime.now().isoformat()}\n\n")

    # 개요
    report.append("## 📊 분석 개요\n\n")
    report.append(f"- 분석된 FreeLang 파일: {len(extractor.all_files)}개\n")
    report.append(f"- 식별된 함수: {len(extractor.function_defs)}개\n")
    report.append(f"- 식별된 모듈: {len(extractor.modules)}개\n")
    report.append(f"- 디자인 패턴: {len(extractor.patterns)}개\n\n")

    # 상위 함수 (반복도)
    report.append("## 🔥 상위 함수 (재사용도 높음)\n\n")
    top_functions = sorted(
        extractor.function_defs.items(),
        key=lambda x: -x[1]
    )[:30]

    for func, count in top_functions:
        report.append(f"- `{func}`: {count}개 프로젝트에서 사용\n")

    report.append("\n")

    # 상위 모듈
    report.append("## 🏗️ 상위 모듈 (반복되는 구조)\n\n")
    top_modules = sorted(
        extractor.modules.items(),
        key=lambda x: -len(x[1])
    )[:20]

    for mod, files in top_modules:
        report.append(f"- `{mod}`: {len(files)}개 파일에서 사용\n")

    report.append("\n")

    # 디자인 패턴 빈도
    report.append("## 🎨 디자인 패턴 분포\n\n")
    top_patterns = sorted(
        extractor.patterns.items(),
        key=lambda x: -x[1]
    )[:15]

    for pattern, count in top_patterns:
        report.append(f"- `{pattern}`: {count}회\n")

    report.append("\n")

    # Standard Core 구성도
    report.append("## 🎯 Standard Core 라이브러리 설계\n\n")
    report.append("### 핵심 모듈 (12개)\n\n")

    total_loc = 0
    for module_name in sorted(core_modules.keys()):
        module = core_modules[module_name]
        loc = estimates.get(module_name, 0)
        total_loc += loc
        freq = module['frequency']

        report.append(f"#### {module_name.upper()}\n")
        report.append(f"- 설명: {module['description']}\n")
        report.append(f"- 함수: {', '.join(module['functions'][:5])}...\n")
        report.append(f"- 실제 사용 빈도: {freq}회\n")
        report.append(f"- 예상 코드량: {loc}줄\n\n")

    report.append(f"### 총 예상 Standard Core 규모\n\n")
    report.append(f"- **코드량**: {total_loc}줄\n")
    report.append(f"- **모듈 수**: {len(core_modules)}개\n")
    report.append(f"- **개발 기간**: 2-3주 (병렬 개발 가능)\n")
    report.append(f"- **테스트**: 각 모듈당 50-100개 테스트 케이스 (총 800-1200개)\n\n")

    # 의존성 그래프
    report.append("## 🔗 모듈 간 의존성\n\n")
    report.append("""
```
result ──────┐
option ──────┤
             ├──> collections (Vec, HashMap, ...)
iter ────────┤       ↓
             ├──> string (문자열 처리)
string ──────┤       ↓
math ────────┤──> io (파일, 콘솔)
             ↓
             ├──> serialization (JSON)
             ├──> crypto (암호화)
             ├──> network (HTTP, TCP)
             └──> async/concurrency
```
""")

    report.append("\n## 📈 마이그레이션 로드맵\n\n")
    report.append("""
### Phase 2-A: Standard Core 설계 (1주)
- [x] 함수/모듈 패턴 분석 (완료)
- [ ] API 설계 (공개 인터페이스)
- [ ] 문서화 및 예제 작성

### Phase 2-B: 핵심 모듈 구현 (2주)
**우선순위 1** (기초):
- collections (Vec, HashMap)
- result/option (에러 처리)
- string (문자열)

**우선순위 2** (활용):
- io (파일 입출력)
- iter (반복자)
- math (수학)

**우선순위 3** (고급):
- async/concurrency (병렬)
- serialization (직렬화)
- crypto (암호화)
- network (통신)

### Phase 2-C: 130개 프로젝트 마이그레이션 (3주)
- 자동 코드 변환 스크립트
- 호환성 테스트
- 점진적 롤아웃

### Phase 3: Z3 SMT 검증 (진행 예정)
""")

    return ''.join(report)


def main():
    root = Path('/data/data/com.termux/files/home')

    # 패턴 추출
    extractor = PatternExtractor(root)
    all_files = extractor.extract_all()

    # Standard Core 설계
    designer = StandardCoreDesigner(extractor)
    core_modules = designer.design()
    estimates = designer.estimate_loc()

    # 리포트 생성
    report = generate_report(extractor, core_modules, estimates)
    with open('phase2_analysis.md', 'w') as f:
        f.write(report)

    # JSON 형식으로도 저장 (자동화용)
    export_data = {
        'files_analyzed': len(all_files),
        'functions_found': len(extractor.function_defs),
        'modules_found': len(extractor.modules),
        'patterns_found': len(extractor.patterns),
        'top_functions': dict(sorted(
            extractor.function_defs.items(),
            key=lambda x: -x[1]
        )[:30]),
        'standard_core': {
            module: {
                'description': info['description'],
                'functions': info['functions'],
                'frequency': info['frequency'],
                'estimated_loc': estimates.get(module, 0),
            }
            for module, info in core_modules.items()
        },
        'total_estimated_loc': sum(estimates.values()),
    }

    with open('phase2_export.json', 'w') as f:
        json.dump(export_data, f, indent=2)

    print("✅ Phase 2 완료!")
    print("📁 생성된 파일:")
    print("  - phase2_analysis.md (상세 분석 리포트)")
    print("  - phase2_export.json (자동화용 데이터)")
    print(f"\n📊 요약:")
    print(f"  - 분석된 파일: {len(all_files)}개")
    print(f"  - Standard Core: {sum(estimates.values())}줄 예상")
    print(f"  - 개발 기간: 2-3주")


if __name__ == '__main__':
    main()
