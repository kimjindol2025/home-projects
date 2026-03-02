#!/usr/bin/env python3
"""
C11 라이브러리 저장소 분석 도구
Gogs API를 통한 자동화된 저장소 분석
"""

import requests
import json
import subprocess
import os
from datetime import datetime
from typing import Dict, List, Tuple, Optional
import csv

class GogsAnalyzer:
    """Gogs 저장소 분석 클래스"""

    def __init__(self, api_url: str, token: str, username: str):
        """
        초기화

        Args:
            api_url: Gogs API 기본 URL (예: https://gogs.dclub.kr/api/v1)
            token: Gogs API 토큰
            username: 사용자명
        """
        self.api_url = api_url
        self.token = token
        self.username = username
        self.headers = {
            "Authorization": f"token {token}",
            "Content-Type": "application/json"
        }
        self.repos = [
            "range-c", "moment-c", "knex-c", "lodash-c", "mime-c",
            "uuid-c", "fs-c", "glob-c", "chalk-c", "async-c",
            "graphql-c", "final-c", "connect-c", "cookie-c",
            "cookie_crypto-c", "cookie-parser-c", "body-c",
            "commander-c", "express-c", "c-libs-ci"
        ]
        self.results = []

    def get_repo_info(self, repo_name: str) -> Optional[Dict]:
        """
        Gogs API에서 저장소 정보 조회

        Args:
            repo_name: 저장소명

        Returns:
            저장소 정보 딕셔너리 또는 None
        """
        try:
            url = f"{self.api_url}/repos/{self.username}/{repo_name}"
            response = requests.get(url, headers=self.headers, timeout=10)

            if response.status_code == 200:
                return response.json()
            else:
                print(f"⚠️  {repo_name}: API 오류 {response.status_code}")
                return None
        except Exception as e:
            print(f"❌ {repo_name}: {str(e)}")
            return None

    def clone_repo(self, repo_name: str, target_dir: str) -> bool:
        """
        Git 저장소 클론

        Args:
            repo_name: 저장소명
            target_dir: 대상 디렉토리

        Returns:
            성공 여부
        """
        try:
            url = f"https://gogs.dclub.kr/{self.username}/{repo_name}.git"
            repo_path = os.path.join(target_dir, repo_name)

            if os.path.exists(repo_path):
                print(f"⏭️  {repo_name}: 이미 존재, 스킵")
                return True

            print(f"📥 {repo_name} 클론 중...")
            subprocess.run(
                ["git", "clone", url, repo_path],
                capture_output=True,
                timeout=30
            )
            return True
        except Exception as e:
            print(f"❌ {repo_name} 클론 실패: {str(e)}")
            return False

    def count_files(self, repo_path: str, extension: str) -> int:
        """
        특정 확장자 파일 개수 계산

        Args:
            repo_path: 저장소 경로
            extension: 파일 확장자 (예: '.c', '.h')

        Returns:
            파일 개수
        """
        try:
            result = subprocess.run(
                f"find {repo_path} -name '*{extension}' | wc -l",
                shell=True,
                capture_output=True,
                text=True,
                timeout=10
            )
            return int(result.stdout.strip())
        except:
            return 0

    def count_lines(self, repo_path: str, extension: str) -> int:
        """
        특정 확장자 파일의 총 라인 수 계산

        Args:
            repo_path: 저장소 경로
            extension: 파일 확장자

        Returns:
            총 라인 수
        """
        try:
            result = subprocess.run(
                f"find {repo_path} -name '*{extension}' -exec wc -l {{}} + | tail -1",
                shell=True,
                capture_output=True,
                text=True,
                timeout=10
            )
            parts = result.stdout.strip().split()
            return int(parts[0]) if parts else 0
        except:
            return 0

    def check_file_exists(self, repo_path: str, filename: str) -> bool:
        """
        파일 존재 여부 확인

        Args:
            repo_path: 저장소 경로
            filename: 파일명

        Returns:
            존재 여부
        """
        return os.path.exists(os.path.join(repo_path, filename))

    def check_directory_exists(self, repo_path: str, dirname: str) -> bool:
        """
        디렉토리 존재 여부 확인

        Args:
            repo_path: 저장소 경로
            dirname: 디렉토리명

        Returns:
            존재 여부
        """
        return os.path.isdir(os.path.join(repo_path, dirname))

    def get_last_commit(self, repo_path: str) -> Tuple[str, str, str]:
        """
        마지막 커밋 정보 조회

        Args:
            repo_path: 저장소 경로

        Returns:
            (해시, 메시지, 날짜) 튜플
        """
        try:
            # 해시
            hash_result = subprocess.run(
                f"cd {repo_path} && git log -1 --format=%H",
                shell=True,
                capture_output=True,
                text=True,
                timeout=10
            )
            commit_hash = hash_result.stdout.strip()[:8]

            # 메시지
            msg_result = subprocess.run(
                f"cd {repo_path} && git log -1 --format=%s",
                shell=True,
                capture_output=True,
                text=True,
                timeout=10
            )
            message = msg_result.stdout.strip()

            # 날짜
            date_result = subprocess.run(
                f"cd {repo_path} && git log -1 --format=%ci",
                shell=True,
                capture_output=True,
                text=True,
                timeout=10
            )
            date = date_result.stdout.strip()[:10]

            return commit_hash, message, date
        except:
            return "N/A", "N/A", "N/A"

    def read_readme(self, repo_path: str) -> str:
        """
        README.md 파일 읽기

        Args:
            repo_path: 저장소 경로

        Returns:
            README 내용
        """
        readme_path = os.path.join(repo_path, "README.md")
        if os.path.exists(readme_path):
            try:
                with open(readme_path, 'r', encoding='utf-8') as f:
                    return f.read()
            except:
                return ""
        return ""

    def evaluate_documentation(self, repo_path: str, readme: str) -> Dict[str, str]:
        """
        문서화 수준 평가

        Args:
            repo_path: 저장소 경로
            readme: README 내용

        Returns:
            평가 딕셔너리
        """
        evaluation = {
            "api_doc": "❌",
            "examples": "❌",
            "install_guide": "❌"
        }

        # API 문서 확인
        if "API" in readme or "api" in readme or "function" in readme:
            evaluation["api_doc"] = "⚠️"
        if os.path.exists(os.path.join(repo_path, "docs", "API.md")):
            evaluation["api_doc"] = "✅"

        # 예제 확인
        if "example" in readme.lower() or "usage" in readme.lower():
            evaluation["examples"] = "⚠️"
        if os.path.isdir(os.path.join(repo_path, "examples")):
            evaluation["examples"] = "✅"

        # 설치 가이드 확인
        if "install" in readme.lower() or "build" in readme.lower():
            evaluation["install_guide"] = "⚠️"
        if "Makefile" in os.listdir(repo_path):
            evaluation["install_guide"] = "✅"

        return evaluation

    def analyze_repo(self, repo_name: str, repo_path: str) -> Dict:
        """
        저장소 상세 분석

        Args:
            repo_name: 저장소명
            repo_path: 저장소 경로

        Returns:
            분석 결과 딕셔너리
        """
        print(f"\n📊 분석 중: {repo_name}")

        # API 정보 조회
        api_info = self.get_repo_info(repo_name)

        # 파일 메트릭
        c_files = self.count_files(repo_path, ".c")
        h_files = self.count_files(repo_path, ".h")
        c_lines = self.count_lines(repo_path, ".c")
        h_lines = self.count_lines(repo_path, ".h")
        total_lines = c_lines + h_lines

        test_lines = 0
        if self.check_directory_exists(repo_path, "tests"):
            test_lines = self.count_lines(os.path.join(repo_path, "tests"), ".c")

        # 파일 구조
        has_readme = self.check_file_exists(repo_path, "README.md")
        has_license = self.check_file_exists(repo_path, "LICENSE")
        has_makefile = self.check_file_exists(repo_path, "Makefile")
        has_src = self.check_directory_exists(repo_path, "src")
        has_include = self.check_directory_exists(repo_path, "include")
        has_tests = self.check_directory_exists(repo_path, "tests")
        has_docs = self.check_directory_exists(repo_path, "docs")
        has_examples = self.check_directory_exists(repo_path, "examples")

        # README 및 문서화 평가
        readme = self.read_readme(repo_path)
        doc_eval = self.evaluate_documentation(repo_path, readme)

        # 커밋 정보
        commit_hash, commit_msg, commit_date = self.get_last_commit(repo_path)

        # 품질 점수 계산
        code_quality = self._calculate_code_quality(repo_path, c_files, total_lines)
        documentation = self._calculate_documentation_score(doc_eval, has_readme)
        testing = self._calculate_testing_score(has_tests, test_lines)
        maintenance = self._calculate_maintenance_score(commit_date)

        average_score = (code_quality + documentation + testing + maintenance) / 4

        # 최종 평가
        if average_score >= 4.5:
            evaluation = "우수"
        elif average_score >= 3.5:
            evaluation = "양호"
        elif average_score >= 2.5:
            evaluation = "개선필요"
        else:
            evaluation = "부실"

        return {
            "name": repo_name,
            "url": f"https://gogs.dclub.kr/{self.username}/{repo_name}",
            "description": api_info.get("description", "N/A") if api_info else "N/A",
            "created_at": api_info.get("created_at", "N/A") if api_info else "N/A",
            "pushed_at": api_info.get("pushed_at", "N/A") if api_info else "N/A",
            "size": api_info.get("size", 0) if api_info else 0,
            "stars": api_info.get("stars_count", 0) if api_info else 0,
            "watchers": api_info.get("watchers_count", 0) if api_info else 0,
            # 파일 구조
            "has_src": "YES" if has_src else "NO",
            "has_include": "YES" if has_include else "NO",
            "has_tests": "YES" if has_tests else "NO",
            "has_docs": "YES" if has_docs else "NO",
            "has_examples": "YES" if has_examples else "NO",
            "has_readme": "YES" if has_readme else "NO",
            "has_license": "YES" if has_license else "NO",
            "has_makefile": "YES" if has_makefile else "NO",
            # 코드 메트릭
            "c_files": c_files,
            "h_files": h_files,
            "total_loc": total_lines,
            "test_loc": test_lines,
            # 문서화
            "api_doc": doc_eval["api_doc"],
            "examples": doc_eval["examples"],
            "install_guide": doc_eval["install_guide"],
            # 커밋
            "last_commit_hash": commit_hash,
            "last_commit_msg": commit_msg,
            "last_commit_date": commit_date,
            # 품질 점수
            "code_quality": code_quality,
            "documentation": documentation,
            "testing": testing,
            "maintenance": maintenance,
            "average_score": round(average_score, 1),
            "evaluation": evaluation
        }

    def _calculate_code_quality(self, repo_path: str, c_files: int, total_lines: int) -> float:
        """코드 품질 점수 계산"""
        if c_files == 0:
            return 1.0
        if c_files < 5:
            return 2.0
        if total_lines < 1000:
            return 3.0
        if total_lines < 5000:
            return 4.0
        return 5.0

    def _calculate_documentation_score(self, doc_eval: Dict, has_readme: bool) -> float:
        """문서화 점수 계산"""
        if not has_readme:
            return 1.0

        score = 2.0
        if doc_eval["api_doc"] == "✅":
            score += 1.0
        elif doc_eval["api_doc"] == "⚠️":
            score += 0.5

        if doc_eval["examples"] == "✅":
            score += 1.0
        elif doc_eval["examples"] == "⚠️":
            score += 0.5

        if doc_eval["install_guide"] == "✅":
            score += 1.0
        elif doc_eval["install_guide"] == "⚠️":
            score += 0.5

        return min(score, 5.0)

    def _calculate_testing_score(self, has_tests: bool, test_lines: int) -> float:
        """테스트 점수 계산"""
        if not has_tests or test_lines == 0:
            return 1.0
        if test_lines < 100:
            return 2.0
        if test_lines < 500:
            return 3.0
        if test_lines < 2000:
            return 4.0
        return 5.0

    def _calculate_maintenance_score(self, last_commit_date: str) -> float:
        """유지보수 점수 계산"""
        if last_commit_date == "N/A":
            return 1.0

        try:
            commit_date = datetime.strptime(last_commit_date, "%Y-%m-%d")
            days_since = (datetime.now() - commit_date).days

            if days_since < 7:
                return 5.0
            elif days_since < 30:
                return 4.0
            elif days_since < 90:
                return 3.0
            elif days_since < 180:
                return 2.0
            else:
                return 1.0
        except:
            return 2.0

    def run_analysis(self, work_dir: str = "./c11_repos"):
        """
        전체 분석 실행

        Args:
            work_dir: 작업 디렉토리
        """
        print("=" * 80)
        print("C11 라이브러리 저장소 분석 시작")
        print("=" * 80)

        os.makedirs(work_dir, exist_ok=True)

        for repo_name in self.repos:
            # 저장소 클론
            if not self.clone_repo(repo_name, work_dir):
                print(f"❌ {repo_name} 클론 실패, 스킵")
                continue

            # 저장소 분석
            repo_path = os.path.join(work_dir, repo_name)
            result = self.analyze_repo(repo_name, repo_path)
            self.results.append(result)

        # 결과 저장
        self.save_results(work_dir)
        self.print_summary()

    def save_results(self, work_dir: str):
        """
        분석 결과를 CSV 및 JSON으로 저장

        Args:
            work_dir: 작업 디렉토리
        """
        # CSV 저장
        csv_path = os.path.join(work_dir, "C11_ANALYSIS_RESULTS.csv")
        with open(csv_path, 'w', newline='', encoding='utf-8') as f:
            writer = csv.DictWriter(f, fieldnames=self.results[0].keys())
            writer.writeheader()
            writer.writerows(self.results)
        print(f"\n✅ CSV 저장됨: {csv_path}")

        # JSON 저장
        json_path = os.path.join(work_dir, "C11_ANALYSIS_RESULTS.json")
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(self.results, f, indent=2, ensure_ascii=False)
        print(f"✅ JSON 저장됨: {json_path}")

    def print_summary(self):
        """분석 결과 요약 출력"""
        print("\n" + "=" * 80)
        print("분석 결과 요약")
        print("=" * 80)

        # 평가별 분류
        excellent = [r for r in self.results if r["evaluation"] == "우수"]
        good = [r for r in self.results if r["evaluation"] == "양호"]
        needs_improvement = [r for r in self.results if r["evaluation"] == "개선필요"]
        poor = [r for r in self.results if r["evaluation"] == "부실"]

        print(f"\n📊 평가 분포:")
        print(f"  우수 (Excellent): {len(excellent)} 개 ({len(excellent)*100/len(self.results):.1f}%)")
        print(f"  양호 (Good): {len(good)} 개 ({len(good)*100/len(self.results):.1f}%)")
        print(f"  개선필요 (Needs Improvement): {len(needs_improvement)} 개 ({len(needs_improvement)*100/len(self.results):.1f}%)")
        print(f"  부실 (Poor): {len(poor)} 개 ({len(poor)*100/len(self.results):.1f}%)")

        # 점수 통계
        scores = [r["average_score"] for r in self.results]
        print(f"\n📈 점수 통계:")
        print(f"  평균: {sum(scores)/len(scores):.2f}")
        print(f"  최고: {max(scores):.2f}")
        print(f"  최저: {min(scores):.2f}")

        # 상위 5개
        print(f"\n🏆 상위 5개 (점수 순):")
        sorted_results = sorted(self.results, key=lambda x: x["average_score"], reverse=True)
        for i, result in enumerate(sorted_results[:5], 1):
            print(f"  {i}. {result['name']}: {result['average_score']}점 ({result['evaluation']})")

        # 하위 5개
        print(f"\n📉 하위 5개 (점수 순):")
        for i, result in enumerate(sorted_results[-5:], 1):
            print(f"  {i}. {result['name']}: {result['average_score']}점 ({result['evaluation']})")


def main():
    """메인 함수"""
    # 설정
    API_URL = "https://gogs.dclub.kr/api/v1"
    TOKEN = "826b3705d8a0602cf89a02327dcee25e991dd630"
    USERNAME = "kim"

    # 분석 실행
    analyzer = GogsAnalyzer(API_URL, TOKEN, USERNAME)
    analyzer.run_analysis()


if __name__ == "__main__":
    main()
