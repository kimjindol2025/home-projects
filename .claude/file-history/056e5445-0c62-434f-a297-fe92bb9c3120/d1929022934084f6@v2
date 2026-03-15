/*
    FreeLiulia File I/O Test Suite

Phase 7: 파일 I/O 함수 20개 테스트
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <unistd.h>

#include "../src/freeliulia_fileio.h"

/* ======================== 테스트 프레임워크 ======================== */

int test_count = 0;
int test_passed = 0;

void test(const char* name, bool condition) {
    test_count++;
    if (condition) {
        test_passed++;
        printf("✓ Test %d: %s\n", test_count, name);
    } else {
        printf("✗ Test %d: %s\n", test_count, name);
    }
}

/* ======================== 메인 함수 ======================== */

int main() {
    printf("\n");
    printf("███████████████████████████████████████████████████████████████\n");
    printf("🚀 Phase 7: 파일 I/O 라이브러리 테스트\n");
    printf("███████████████████████████████████████████████████████████████\n");

/* ======================== Test Suite 1: 파일 생성 및 쓰기 ======================== */

    printf("\n");
    printf("📌 Test Suite 1: 파일 생성 및 쓰기 (5 tests)\n");
printf("─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "\n");

/* Test 1.1: 파일 생성 */
{
    파일핸들 f = 파일열기("/tmp/test_fl_1.txt", "w");
    test("File creation with mode 'w'", f != NULL);
    if (f != NULL) 파일닫기(f);
}

/* Test 1.2: 파일에 한 줄 쓰기 */
{
    파일핸들 f = 파일열기("/tmp/test_fl_2.txt", "w");
    bool result = false;
    if (f != NULL) {
        int64_t ret = 줄쓰기(f, "Hello World");
        result = (ret >= 0);
        파일닫기(f);
    }
    test("Write single line to file", result);
}

/* Test 1.3: 여러 줄 쓰기 */
{
    파일핸들 f = 파일열기("/tmp/test_fl_3.txt", "w");
    bool result = true;
    if (f != NULL) {
        줄쓰기(f, "Line 1");
        줄쓰기(f, "Line 2");
        줄쓰기(f, "Line 3");
        파일닫기(f);
    } else {
        result = false;
    }
    test("Write multiple lines to file", result);
}

/* Test 1.4: 빈 파일 생성 */
{
    파일핸들 f = 파일열기("/tmp/test_fl_4.txt", "w");
    bool result = false;
    if (f != NULL) {
        파일닫기(f);
        result = 파일존재("/tmp/test_fl_4.txt");
    }
    test("Create empty file", result);
}

/* Test 1.5: 파일 크기 확인 */
{
    파일핸들 f = 파일열기("/tmp/test_fl_5.txt", "w");
    if (f != NULL) {
        줄쓰기(f, "Test");
        파일닫기(f);
    }
    int64_t size = 파일크기("/tmp/test_fl_5.txt");
    test("Get file size", size > 0);
}

/* ======================== Test Suite 2: 파일 읽기 ======================== */

printf("\n");
printf("📌 Test Suite 2: 파일 읽기 (5 tests)\n");
printf("─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "\n");

/* Test 2.1: 파일 열기 읽기 모드 */
{
    /* 먼저 테스트 파일 생성 */
    파일핸들 fw = 파일열기("/tmp/test_fl_read_1.txt", "w");
    줄쓰기(fw, "Test content");
    파일닫기(fw);

    파일핸들 f = 파일열기("/tmp/test_fl_read_1.txt", "r");
    test("Open file in read mode", f != NULL);
    if (f != NULL) 파일닫기(f);
}

/* Test 2.2: 파일에서 한 줄 읽기 */
{
    파일핸들 fw = 파일열기("/tmp/test_fl_read_2.txt", "w");
    줄쓰기(fw, "Hello from file");
    파일닫기(fw);

    파일핸들 f = 파일열기("/tmp/test_fl_read_2.txt", "r");
    bool result = false;
    if (f != NULL) {
        char* line = 줄읽기(f);
        result = (line != NULL && strlen(line) > 0);
        if (line != NULL) free(line);
        파일닫기(f);
    }
    test("Read single line from file", result);
}

/* Test 2.3: 빈 파일 읽기 */
{
    파일핸들 fw = 파일열기("/tmp/test_fl_read_3.txt", "w");
    파일닫기(fw);

    파일핸들 f = 파일열기("/tmp/test_fl_read_3.txt", "r");
    bool result = false;
    if (f != NULL) {
        char* line = 줄읽기(f);
        result = (line == NULL);
        if (line != NULL) free(line);
        파일닫기(f);
    }
    test("Read from empty file returns NULL", result);
}

/* Test 2.4: 여러 줄 읽기 */
{
    파일핸들 fw = 파일열기("/tmp/test_fl_read_4.txt", "w");
    줄쓰기(fw, "Line 1");
    줄쓰기(fw, "Line 2");
    파일닫기(fw);

    파일핸들 f = 파일열기("/tmp/test_fl_read_4.txt", "r");
    bool result = false;
    if (f != NULL) {
        char* line1 = 줄읽기(f);
        char* line2 = 줄읽기(f);
        result = (line1 != NULL && line2 != NULL);
        if (line1 != NULL) free(line1);
        if (line2 != NULL) free(line2);
        파일닫기(f);
    }
    test("Read multiple lines from file", result);
}

/* Test 2.5: 전체 파일 읽기 */
{
    파일핸들 fw = 파일열기("/tmp/test_fl_read_5.txt", "w");
    줄쓰기(fw, "Complete");
    파일닫기(fw);

    char* content = 파일전체읽기("/tmp/test_fl_read_5.txt");
    bool result = (content != NULL && strlen(content) > 0);
    if (content != NULL) free(content);
    test("Read entire file into memory", result);
}

/* ======================== Test Suite 3: 파일 존재 및 삭제 ======================== */

printf("\n");
printf("📌 Test Suite 3: 파일 존재 및 삭제 (5 tests)\n");
printf("─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "\n");

/* Test 3.1: 파일 존재 확인 */
{
    파일핸들 f = 파일열기("/tmp/test_fl_exist_1.txt", "w");
    파일닫기(f);
    bool result = 파일존재("/tmp/test_fl_exist_1.txt");
    test("File exists after creation", result);
}

/* Test 3.2: 없는 파일 확인 */
{
    bool result = !파일존재("/tmp/nonexistent_file_xyz_12345.txt");
    test("Nonexistent file returns false", result);
}

/* Test 3.3: 파일 삭제 */
{
    파일핸들 f = 파일열기("/tmp/test_fl_delete_1.txt", "w");
    파일닫기(f);

    int64_t ret = 파일삭제("/tmp/test_fl_delete_1.txt");
    bool result = (ret == 0 && !파일존재("/tmp/test_fl_delete_1.txt"));
    test("Delete existing file", result);
}

/* Test 3.4: 없는 파일 삭제 시도 */
{
    int64_t ret = 파일삭제("/tmp/nonexistent_file_xyz_67890.txt");
    bool result = (ret != 0);  /* 실패해야 함 */
    test("Delete nonexistent file returns error", result);
}

/* Test 3.5: 파일 크기 0인 파일 */
{
    파일핸들 f = 파일열기("/tmp/test_fl_size_0.txt", "w");
    파일닫기(f);
    int64_t size = 파일크기("/tmp/test_fl_size_0.txt");
    test("Empty file has size 0", size == 0);
}

/* ======================== Test Suite 4: 파일 모드 ======================== */

printf("\n");
printf("📌 Test Suite 4: 파일 모드 (5 tests)\n");
printf("─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "─" "\n");

/* Test 4.1: Append 모드 */
{
    파일핸들 fw = 파일열기("/tmp/test_fl_append.txt", "w");
    줄쓰기(fw, "First");
    파일닫기(fw);

    파일핸들 fa = 파일열기("/tmp/test_fl_append.txt", "a");
    줄쓰기(fa, "Second");
    파일닫기(fa);

    int64_t size = 파일크기("/tmp/test_fl_append.txt");
    test("Append mode adds to file", size > 10);  /* First + Second + 2 newlines */
}

/* Test 4.2: 쓰기 모드는 덮어쓰기 */
{
    파일핸들 fw1 = 파일열기("/tmp/test_fl_overwrite.txt", "w");
    줄쓰기(fw1, "Original content");
    파일닫기(fw1);

    파일핸들 fw2 = 파일열기("/tmp/test_fl_overwrite.txt", "w");
    줄쓰기(fw2, "New");
    파일닫기(fw2);

    int64_t size = 파일크기("/tmp/test_fl_overwrite.txt");
    test("Write mode overwrites file", size < 20);  /* "New" is shorter than "Original..." */
}

/* Test 4.3: NULL 경로 처리 */
{
    파일핸들 f = 파일열기(NULL, "r");
    test("NULL path returns NULL handle", f == NULL);
}

/* Test 4.4: NULL 모드 처리 */
{
    파일핸들 f = 파일열기("/tmp/test_fl_null_mode.txt", NULL);
    test("NULL mode returns NULL handle", f == NULL);
}

/* Test 4.5: NULL 핸들 닫기 */
{
    int64_t ret = 파일닫기(NULL);
    test("Closing NULL handle returns -1", ret == -1);
}

/* ======================== 테스트 결과 요약 ======================== */

    printf("\n");
    printf("=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "\n");
    printf("✨ File I/O Test Summary\n");
    printf("=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "\n");
    printf("✅ 통과: %d / %d\n", test_passed, test_count);

    if (test_passed == test_count) {
        printf("🎉 모든 파일 I/O 테스트 통과!\n");
    } else {
        printf("⚠️  일부 테스트 실패: %d/%d\n", test_count - test_passed, test_count);
    }

    printf("=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "=" "\n");

    /* 정리: 테스트 파일 삭제 */
    system("rm -f /tmp/test_fl_*.txt");

    return (test_passed == test_count) ? 0 : 1;
}
