#ifndef FREELIULIA_FILEIO_H
#define FREELIULIA_FILEIO_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include <unistd.h>

/* ======================== 파일 I/O 함수 ======================== */

/* 파일 핸들 구조 (불투명 타입) */
typedef FILE* 파일핸들;

/* 파일 열기: 경로, 모드("r", "w", "a") */
static inline 파일핸들 파일열기(const char* 경로, const char* 모드) {
    if (경로 == NULL || 모드 == NULL) {
        return NULL;
    }
    파일핸들 handle = fopen(경로, 모드);
    return handle;
}

/* 파일 닫기 */
static inline int64_t 파일닫기(파일핸들 handle) {
    if (handle == NULL) {
        return -1;
    }
    int result = fclose(handle);
    return (int64_t)result;
}

/* 파일에서 한 줄 읽기 (최대 1024자) */
static inline char* 줄읽기(파일핸들 handle) {
    if (handle == NULL) {
        return NULL;
    }
    char* buffer = (char*)malloc(1024);
    if (buffer == NULL) {
        return NULL;
    }

    char* result = fgets(buffer, 1024, handle);
    if (result == NULL) {
        free(buffer);
        return NULL;
    }

    return buffer;
}

/* 파일에 한 줄 쓰기 */
static inline int64_t 줄쓰기(파일핸들 handle, const char* 내용) {
    if (handle == NULL || 내용 == NULL) {
        return -1;
    }

    int result = fputs(내용, handle);
    if (result < 0) {
        return -1;
    }

    /* 줄 끝 추가 */
    int nl_result = fputc('\n', handle);
    if (nl_result < 0) {
        return -1;
    }

    return (int64_t)result;
}

/* 파일이 존재하는지 확인 */
static inline bool 파일존재(const char* 경로) {
    if (경로 == NULL) {
        return false;
    }

    /* access() 사용: F_OK = 파일 존재 확인 */
    return access(경로, F_OK) == 0;
}

/* 파일 삭제 */
static inline int64_t 파일삭제(const char* 경로) {
    if (경로 == NULL) {
        return -1;
    }

    int result = remove(경로);
    return (int64_t)result;
}

/* 파일 크기 반환 (바이트) */
static inline int64_t 파일크기(const char* 경로) {
    if (경로 == NULL) {
        return -1;
    }

    FILE* f = fopen(경로, "rb");
    if (f == NULL) {
        return -1;
    }

    fseek(f, 0, SEEK_END);
    long size = ftell(f);
    fclose(f);

    return (size < 0) ? -1 : (int64_t)size;
}

/* 파일의 모든 내용 읽기 (메모리 할당) */
static inline char* 파일전체읽기(const char* 경로) {
    if (경로 == NULL) {
        return NULL;
    }

    FILE* f = fopen(경로, "rb");
    if (f == NULL) {
        return NULL;
    }

    fseek(f, 0, SEEK_END);
    long size = ftell(f);
    fseek(f, 0, SEEK_SET);

    char* buffer = (char*)malloc(size + 1);
    if (buffer == NULL) {
        fclose(f);
        return NULL;
    }

    size_t read = fread(buffer, 1, size, f);
    buffer[read] = '\0';
    fclose(f);

    return buffer;
}

#endif /* FREELIULIA_FILEIO_H */
