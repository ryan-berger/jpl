; ModuleID = 'runtime.c'
source_filename = "runtime.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-pc-linux-gnu"

%struct.anon = type { i64, [256 x i8] }
%struct._IO_FILE = type { i32, i8*, i8*, i8*, i8*, i8*, i8*, i8*, i8*, i8*, i8*, i8*, %struct._IO_marker*, %struct._IO_FILE*, i32, i32, i64, i16, i8, [1 x i8], i8*, i64, %struct._IO_codecvt*, %struct._IO_wide_data*, %struct._IO_FILE*, i8*, i64, i32, [20 x i8] }
%struct._IO_marker = type opaque
%struct._IO_codecvt = type opaque
%struct._IO_wide_data = type opaque
%struct.__va_list_tag = type { i32, i32, i8*, i8* }
%struct.pict = type { i64, i64, double* }

@PB = dso_local constant i32 1, align 4
@MAX = dso_local constant i32 256, align 4
@mem = common dso_local global %struct.anon zeroinitializer, align 8
@stderr = external dso_local global %struct._IO_FILE*, align 8
@.str = private unnamed_addr constant [46 x i8] c"[builtin show] Type too complex, cannot parse\00", align 1
@.str.1 = private unnamed_addr constant [39 x i8] c"[builtin show] printf failed, aborting\00", align 1
@.str.2 = private unnamed_addr constant [50 x i8] c"[builtin show] Could not parse boolean type in %s\00", align 1
@.str.3 = private unnamed_addr constant [50 x i8] c"[builtin show] Could not parse integer type in %s\00", align 1
@.str.4 = private unnamed_addr constant [51 x i8] c"[builtin show] Could not parse floating type in %s\00", align 1
@.str.5 = private unnamed_addr constant [60 x i8] c"[builtin show] Tuples with %zu or more fields not supported\00", align 1
@.str.6 = private unnamed_addr constant [48 x i8] c"[builtin show] Could not parse tuple type in %s\00", align 1
@.str.7 = private unnamed_addr constant [42 x i8] c"[builtin show] Could not parse type in %s\00", align 1
@.str.8 = private unnamed_addr constant [45 x i8] c"[builtin show] Ranks above 255 not supported\00", align 1
@.str.9 = private unnamed_addr constant [59 x i8] c"[builtin show] Overflow when computing total size of array\00", align 1
@.str.10 = private unnamed_addr constant [2 x i8] c"[\00", align 1
@.str.11 = private unnamed_addr constant [3 x i8] c", \00", align 1
@.str.12 = private unnamed_addr constant [2 x i8] c";\00", align 1
@.str.13 = private unnamed_addr constant [2 x i8] c" \00", align 1
@.str.14 = private unnamed_addr constant [2 x i8] c"]\00", align 1
@.str.15 = private unnamed_addr constant [5 x i8] c"true\00", align 1
@.str.16 = private unnamed_addr constant [6 x i8] c"false\00", align 1
@.str.17 = private unnamed_addr constant [5 x i8] c"%lld\00", align 1
@.str.18 = private unnamed_addr constant [3 x i8] c"%f\00", align 1
@.str.19 = private unnamed_addr constant [2 x i8] c"{\00", align 1
@.str.20 = private unnamed_addr constant [2 x i8] c"}\00", align 1
@.str.21 = private unnamed_addr constant [45 x i8] c"[builtin show] Type string is too long in %s\00", align 1
@.str.22 = private unnamed_addr constant [46 x i8] c"[builtin show] Cannot parse type string in %s\00", align 1
@.str.23 = private unnamed_addr constant [11 x i8] c"[abort] %s\00", align 1
@.str.24 = private unnamed_addr constant [3 x i8] c"%s\00", align 1
@.str.25 = private unnamed_addr constant [15 x i8] c"malloc failed\0A\00", align 1
@.str.26 = private unnamed_addr constant [30 x i8] c"Blur radius must be positive\0A\00", align 1
@.str.27 = private unnamed_addr constant [32 x i8] c"Must resize to a positive size\0A\00", align 1
@.str.28 = private unnamed_addr constant [19 x i8] c"reading image: %s\0A\00", align 1
@.str.29 = private unnamed_addr constant [18 x i8] c"writing image: %s\00", align 1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local zeroext i8 @getmem(i64 %0) #0 {
  %2 = alloca i64, align 8
  store i64 %0, i64* %2, align 8
  %3 = load i64, i64* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 0), align 8
  %4 = load i64, i64* %2, align 8
  %5 = add i64 %3, %4
  %6 = icmp ult i64 %5, 256
  br i1 %6, label %7, label %15

7:                                                ; preds = %1
  %8 = load i64, i64* %2, align 8
  %9 = load i64, i64* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 0), align 8
  %10 = add i64 %9, %8
  store i64 %10, i64* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 0), align 8
  %11 = load i64, i64* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 0), align 8
  %12 = load i64, i64* %2, align 8
  %13 = sub i64 %11, %12
  %14 = trunc i64 %13 to i8
  ret i8 %14

15:                                               ; preds = %1
  %16 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %17 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %16, i8* getelementptr inbounds ([46 x i8], [46 x i8]* @.str, i64 0, i64 0))
  call void @exit(i32 127) #7
  unreachable
}

declare dso_local i32 @fprintf(%struct._IO_FILE*, i8*, ...) #1

; Function Attrs: noreturn nounwind
declare dso_local void @exit(i32) #2

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @tprintf(i8* %0, ...) #0 {
  %2 = alloca i8*, align 8
  %3 = alloca i32, align 4
  %4 = alloca [1 x %struct.__va_list_tag], align 16
  store i8* %0, i8** %2, align 8
  store i32 0, i32* %3, align 4
  %5 = getelementptr inbounds [1 x %struct.__va_list_tag], [1 x %struct.__va_list_tag]* %4, i64 0, i64 0
  %6 = bitcast %struct.__va_list_tag* %5 to i8*
  call void @llvm.va_start(i8* %6)
  %7 = load i8*, i8** %2, align 8
  %8 = getelementptr inbounds [1 x %struct.__va_list_tag], [1 x %struct.__va_list_tag]* %4, i64 0, i64 0
  %9 = call i32 @vprintf(i8* %7, %struct.__va_list_tag* %8)
  %10 = icmp slt i32 %9, 0
  br i1 %10, label %11, label %12

11:                                               ; preds = %1
  store i32 1, i32* %3, align 4
  br label %12

12:                                               ; preds = %11, %1
  %13 = getelementptr inbounds [1 x %struct.__va_list_tag], [1 x %struct.__va_list_tag]* %4, i64 0, i64 0
  %14 = bitcast %struct.__va_list_tag* %13 to i8*
  call void @llvm.va_end(i8* %14)
  %15 = load i32, i32* %3, align 4
  %16 = icmp ne i32 %15, 0
  br i1 %16, label %17, label %20

17:                                               ; preds = %12
  %18 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %19 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %18, i8* getelementptr inbounds ([39 x i8], [39 x i8]* @.str.1, i64 0, i64 0))
  call void @exit(i32 127) #7
  unreachable

20:                                               ; preds = %12
  ret void
}

; Function Attrs: nounwind
declare void @llvm.va_start(i8*) #3

declare dso_local i32 @vprintf(i8*, %struct.__va_list_tag*) #1

; Function Attrs: nounwind
declare void @llvm.va_end(i8*) #3

; Function Attrs: noinline nounwind optnone uwtable
define dso_local zeroext i8 @parse_bool_type(i8* %0, i8** %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8**, align 8
  %5 = alloca i8*, align 8
  %6 = alloca i8, align 1
  store i8* %0, i8** %3, align 8
  store i8** %1, i8*** %4, align 8
  %7 = load i8*, i8** %3, align 8
  store i8* %7, i8** %5, align 8
  %8 = load i8*, i8** %3, align 8
  %9 = getelementptr inbounds i8, i8* %8, i32 1
  store i8* %9, i8** %3, align 8
  %10 = load i8, i8* %8, align 1
  %11 = sext i8 %10 to i32
  %12 = icmp ne i32 %11, 98
  br i1 %12, label %13, label %14

13:                                               ; preds = %2
  br label %43

14:                                               ; preds = %2
  %15 = load i8*, i8** %3, align 8
  %16 = getelementptr inbounds i8, i8* %15, i32 1
  store i8* %16, i8** %3, align 8
  %17 = load i8, i8* %15, align 1
  %18 = sext i8 %17 to i32
  %19 = icmp ne i32 %18, 111
  br i1 %19, label %20, label %21

20:                                               ; preds = %14
  br label %43

21:                                               ; preds = %14
  %22 = load i8*, i8** %3, align 8
  %23 = getelementptr inbounds i8, i8* %22, i32 1
  store i8* %23, i8** %3, align 8
  %24 = load i8, i8* %22, align 1
  %25 = sext i8 %24 to i32
  %26 = icmp ne i32 %25, 111
  br i1 %26, label %27, label %28

27:                                               ; preds = %21
  br label %43

28:                                               ; preds = %21
  %29 = load i8*, i8** %3, align 8
  %30 = getelementptr inbounds i8, i8* %29, i32 1
  store i8* %30, i8** %3, align 8
  %31 = load i8, i8* %29, align 1
  %32 = sext i8 %31 to i32
  %33 = icmp ne i32 %32, 108
  br i1 %33, label %34, label %35

34:                                               ; preds = %28
  br label %43

35:                                               ; preds = %28
  %36 = call zeroext i8 @getmem(i64 1)
  store i8 %36, i8* %6, align 1
  %37 = load i8, i8* %6, align 1
  %38 = zext i8 %37 to i64
  %39 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %38
  store i8 -5, i8* %39, align 1
  %40 = load i8*, i8** %3, align 8
  %41 = load i8**, i8*** %4, align 8
  store i8* %40, i8** %41, align 8
  %42 = load i8, i8* %6, align 1
  ret i8 %42

43:                                               ; preds = %34, %27, %20, %13
  %44 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %45 = load i8*, i8** %5, align 8
  %46 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %44, i8* getelementptr inbounds ([50 x i8], [50 x i8]* @.str.2, i64 0, i64 0), i8* %45)
  call void @exit(i32 127) #7
  unreachable
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local zeroext i8 @parse_int_type(i8* %0, i8** %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8**, align 8
  %5 = alloca i8*, align 8
  %6 = alloca i8, align 1
  store i8* %0, i8** %3, align 8
  store i8** %1, i8*** %4, align 8
  %7 = load i8*, i8** %3, align 8
  store i8* %7, i8** %5, align 8
  %8 = load i8*, i8** %3, align 8
  %9 = getelementptr inbounds i8, i8* %8, i32 1
  store i8* %9, i8** %3, align 8
  %10 = load i8, i8* %8, align 1
  %11 = sext i8 %10 to i32
  %12 = icmp ne i32 %11, 105
  br i1 %12, label %13, label %14

13:                                               ; preds = %2
  br label %36

14:                                               ; preds = %2
  %15 = load i8*, i8** %3, align 8
  %16 = getelementptr inbounds i8, i8* %15, i32 1
  store i8* %16, i8** %3, align 8
  %17 = load i8, i8* %15, align 1
  %18 = sext i8 %17 to i32
  %19 = icmp ne i32 %18, 110
  br i1 %19, label %20, label %21

20:                                               ; preds = %14
  br label %36

21:                                               ; preds = %14
  %22 = load i8*, i8** %3, align 8
  %23 = getelementptr inbounds i8, i8* %22, i32 1
  store i8* %23, i8** %3, align 8
  %24 = load i8, i8* %22, align 1
  %25 = sext i8 %24 to i32
  %26 = icmp ne i32 %25, 116
  br i1 %26, label %27, label %28

27:                                               ; preds = %21
  br label %36

28:                                               ; preds = %21
  %29 = call zeroext i8 @getmem(i64 1)
  store i8 %29, i8* %6, align 1
  %30 = load i8, i8* %6, align 1
  %31 = zext i8 %30 to i64
  %32 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %31
  store i8 -4, i8* %32, align 1
  %33 = load i8*, i8** %3, align 8
  %34 = load i8**, i8*** %4, align 8
  store i8* %33, i8** %34, align 8
  %35 = load i8, i8* %6, align 1
  ret i8 %35

36:                                               ; preds = %27, %20, %13
  %37 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %38 = load i8*, i8** %5, align 8
  %39 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %37, i8* getelementptr inbounds ([50 x i8], [50 x i8]* @.str.3, i64 0, i64 0), i8* %38)
  call void @exit(i32 127) #7
  unreachable
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local zeroext i8 @parse_float_type(i8* %0, i8** %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8**, align 8
  %5 = alloca i8*, align 8
  %6 = alloca i8, align 1
  store i8* %0, i8** %3, align 8
  store i8** %1, i8*** %4, align 8
  %7 = load i8*, i8** %3, align 8
  store i8* %7, i8** %5, align 8
  %8 = load i8*, i8** %3, align 8
  %9 = getelementptr inbounds i8, i8* %8, i32 1
  store i8* %9, i8** %3, align 8
  %10 = load i8, i8* %8, align 1
  %11 = sext i8 %10 to i32
  %12 = icmp ne i32 %11, 102
  br i1 %12, label %13, label %14

13:                                               ; preds = %2
  br label %50

14:                                               ; preds = %2
  %15 = load i8*, i8** %3, align 8
  %16 = getelementptr inbounds i8, i8* %15, i32 1
  store i8* %16, i8** %3, align 8
  %17 = load i8, i8* %15, align 1
  %18 = sext i8 %17 to i32
  %19 = icmp ne i32 %18, 108
  br i1 %19, label %20, label %21

20:                                               ; preds = %14
  br label %50

21:                                               ; preds = %14
  %22 = load i8*, i8** %3, align 8
  %23 = getelementptr inbounds i8, i8* %22, i32 1
  store i8* %23, i8** %3, align 8
  %24 = load i8, i8* %22, align 1
  %25 = sext i8 %24 to i32
  %26 = icmp ne i32 %25, 111
  br i1 %26, label %27, label %28

27:                                               ; preds = %21
  br label %50

28:                                               ; preds = %21
  %29 = load i8*, i8** %3, align 8
  %30 = getelementptr inbounds i8, i8* %29, i32 1
  store i8* %30, i8** %3, align 8
  %31 = load i8, i8* %29, align 1
  %32 = sext i8 %31 to i32
  %33 = icmp ne i32 %32, 97
  br i1 %33, label %34, label %35

34:                                               ; preds = %28
  br label %50

35:                                               ; preds = %28
  %36 = load i8*, i8** %3, align 8
  %37 = getelementptr inbounds i8, i8* %36, i32 1
  store i8* %37, i8** %3, align 8
  %38 = load i8, i8* %36, align 1
  %39 = sext i8 %38 to i32
  %40 = icmp ne i32 %39, 116
  br i1 %40, label %41, label %42

41:                                               ; preds = %35
  br label %50

42:                                               ; preds = %35
  %43 = call zeroext i8 @getmem(i64 1)
  store i8 %43, i8* %6, align 1
  %44 = load i8, i8* %6, align 1
  %45 = zext i8 %44 to i64
  %46 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %45
  store i8 -3, i8* %46, align 1
  %47 = load i8*, i8** %3, align 8
  %48 = load i8**, i8*** %4, align 8
  store i8* %47, i8** %48, align 8
  %49 = load i8, i8* %6, align 1
  ret i8 %49

50:                                               ; preds = %41, %34, %27, %20, %13
  %51 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %52 = load i8*, i8** %5, align 8
  %53 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %51, i8* getelementptr inbounds ([51 x i8], [51 x i8]* @.str.4, i64 0, i64 0), i8* %52)
  call void @exit(i32 127) #7
  unreachable
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local zeroext i8 @parse_tuple_type(i8* %0, i8** %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8**, align 8
  %5 = alloca [256 x i8], align 16
  %6 = alloca i8*, align 8
  %7 = alloca i64, align 8
  %8 = alloca i8, align 1
  %9 = alloca i8, align 1
  %10 = alloca i32, align 4
  store i8* %0, i8** %3, align 8
  store i8** %1, i8*** %4, align 8
  %11 = load i8*, i8** %3, align 8
  store i8* %11, i8** %6, align 8
  %12 = load i8*, i8** %3, align 8
  %13 = getelementptr inbounds i8, i8* %12, i32 1
  store i8* %13, i8** %3, align 8
  %14 = load i8, i8* %12, align 1
  %15 = sext i8 %14 to i32
  %16 = icmp ne i32 %15, 123
  br i1 %16, label %17, label %18

17:                                               ; preds = %2
  br label %115

18:                                               ; preds = %2
  br label %19

19:                                               ; preds = %24, %18
  %20 = load i8*, i8** %3, align 8
  %21 = load i8, i8* %20, align 1
  %22 = sext i8 %21 to i32
  %23 = icmp eq i32 %22, 32
  br i1 %23, label %24, label %27

24:                                               ; preds = %19
  %25 = load i8*, i8** %3, align 8
  %26 = getelementptr inbounds i8, i8* %25, i32 1
  store i8* %26, i8** %3, align 8
  br label %19

27:                                               ; preds = %19
  store i64 0, i64* %7, align 8
  %28 = load i8*, i8** %3, align 8
  %29 = load i8, i8* %28, align 1
  %30 = sext i8 %29 to i32
  %31 = icmp ne i32 %30, 125
  br i1 %31, label %32, label %82

32:                                               ; preds = %27
  %33 = load i8*, i8** %3, align 8
  %34 = call zeroext i8 @parse_type(i8* %33, i8** %3)
  %35 = load i64, i64* %7, align 8
  %36 = add i64 %35, 1
  store i64 %36, i64* %7, align 8
  %37 = getelementptr inbounds [256 x i8], [256 x i8]* %5, i64 0, i64 %35
  store i8 %34, i8* %37, align 1
  br label %38

38:                                               ; preds = %43, %32
  %39 = load i8*, i8** %3, align 8
  %40 = load i8, i8* %39, align 1
  %41 = sext i8 %40 to i32
  %42 = icmp eq i32 %41, 32
  br i1 %42, label %43, label %46

43:                                               ; preds = %38
  %44 = load i8*, i8** %3, align 8
  %45 = getelementptr inbounds i8, i8* %44, i32 1
  store i8* %45, i8** %3, align 8
  br label %38

46:                                               ; preds = %38
  br label %47

47:                                               ; preds = %80, %46
  %48 = load i8*, i8** %3, align 8
  %49 = load i8, i8* %48, align 1
  %50 = sext i8 %49 to i32
  %51 = icmp ne i32 %50, 125
  br i1 %51, label %52, label %81

52:                                               ; preds = %47
  %53 = load i8*, i8** %3, align 8
  %54 = getelementptr inbounds i8, i8* %53, i32 1
  store i8* %54, i8** %3, align 8
  %55 = load i8, i8* %53, align 1
  %56 = sext i8 %55 to i32
  %57 = icmp ne i32 %56, 44
  br i1 %57, label %58, label %59

58:                                               ; preds = %52
  br label %115

59:                                               ; preds = %52
  br label %60

60:                                               ; preds = %65, %59
  %61 = load i8*, i8** %3, align 8
  %62 = load i8, i8* %61, align 1
  %63 = sext i8 %62 to i32
  %64 = icmp eq i32 %63, 32
  br i1 %64, label %65, label %68

65:                                               ; preds = %60
  %66 = load i8*, i8** %3, align 8
  %67 = getelementptr inbounds i8, i8* %66, i32 1
  store i8* %67, i8** %3, align 8
  br label %60

68:                                               ; preds = %60
  %69 = load i8*, i8** %3, align 8
  %70 = call zeroext i8 @parse_type(i8* %69, i8** %3)
  %71 = load i64, i64* %7, align 8
  %72 = add i64 %71, 1
  store i64 %72, i64* %7, align 8
  %73 = getelementptr inbounds [256 x i8], [256 x i8]* %5, i64 0, i64 %71
  store i8 %70, i8* %73, align 1
  %74 = load i64, i64* %7, align 8
  %75 = icmp uge i64 %74, 251
  br i1 %75, label %76, label %80

76:                                               ; preds = %68
  %77 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %78 = load i64, i64* %7, align 8
  %79 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %77, i8* getelementptr inbounds ([60 x i8], [60 x i8]* @.str.5, i64 0, i64 0), i64 %78)
  call void @exit(i32 127) #7
  unreachable

80:                                               ; preds = %68
  br label %47

81:                                               ; preds = %47
  br label %82

82:                                               ; preds = %81, %27
  %83 = load i8*, i8** %3, align 8
  %84 = getelementptr inbounds i8, i8* %83, i32 1
  store i8* %84, i8** %3, align 8
  %85 = load i64, i64* %7, align 8
  %86 = add i64 1, %85
  %87 = call zeroext i8 @getmem(i64 %86)
  store i8 %87, i8* %8, align 1
  store i8 %87, i8* %9, align 1
  %88 = load i64, i64* %7, align 8
  %89 = trunc i64 %88 to i8
  %90 = load i8, i8* %8, align 1
  %91 = add i8 %90, 1
  store i8 %91, i8* %8, align 1
  %92 = zext i8 %90 to i64
  %93 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %92
  store i8 %89, i8* %93, align 1
  store i32 0, i32* %10, align 4
  br label %94

94:                                               ; preds = %108, %82
  %95 = load i32, i32* %10, align 4
  %96 = sext i32 %95 to i64
  %97 = load i64, i64* %7, align 8
  %98 = icmp ult i64 %96, %97
  br i1 %98, label %99, label %111

99:                                               ; preds = %94
  %100 = load i32, i32* %10, align 4
  %101 = sext i32 %100 to i64
  %102 = getelementptr inbounds [256 x i8], [256 x i8]* %5, i64 0, i64 %101
  %103 = load i8, i8* %102, align 1
  %104 = load i8, i8* %8, align 1
  %105 = add i8 %104, 1
  store i8 %105, i8* %8, align 1
  %106 = zext i8 %104 to i64
  %107 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %106
  store i8 %103, i8* %107, align 1
  br label %108

108:                                              ; preds = %99
  %109 = load i32, i32* %10, align 4
  %110 = add nsw i32 %109, 1
  store i32 %110, i32* %10, align 4
  br label %94

111:                                              ; preds = %94
  %112 = load i8*, i8** %3, align 8
  %113 = load i8**, i8*** %4, align 8
  store i8* %112, i8** %113, align 8
  %114 = load i8, i8* %9, align 1
  ret i8 %114

115:                                              ; preds = %58, %17
  %116 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %117 = load i8*, i8** %6, align 8
  %118 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %116, i8* getelementptr inbounds ([48 x i8], [48 x i8]* @.str.6, i64 0, i64 0), i8* %117)
  call void @exit(i32 127) #7
  unreachable
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local zeroext i8 @parse_type(i8* %0, i8** %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8**, align 8
  %5 = alloca i8*, align 8
  %6 = alloca i8, align 1
  %7 = alloca i8, align 1
  %8 = alloca i8, align 1
  %9 = alloca i8, align 1
  store i8* %0, i8** %3, align 8
  store i8** %1, i8*** %4, align 8
  %10 = load i8*, i8** %3, align 8
  store i8* %10, i8** %5, align 8
  br label %11

11:                                               ; preds = %16, %2
  %12 = load i8*, i8** %3, align 8
  %13 = load i8, i8* %12, align 1
  %14 = sext i8 %13 to i32
  %15 = icmp eq i32 %14, 32
  br i1 %15, label %16, label %19

16:                                               ; preds = %11
  %17 = load i8*, i8** %3, align 8
  %18 = getelementptr inbounds i8, i8* %17, i32 1
  store i8* %18, i8** %3, align 8
  br label %11

19:                                               ; preds = %11
  store i8 0, i8* %6, align 1
  %20 = load i8*, i8** %3, align 8
  %21 = load i8, i8* %20, align 1
  %22 = sext i8 %21 to i32
  switch i32 %22, label %35 [
    i32 123, label %23
    i32 105, label %26
    i32 102, label %29
    i32 98, label %32
  ]

23:                                               ; preds = %19
  %24 = load i8*, i8** %3, align 8
  %25 = call zeroext i8 @parse_tuple_type(i8* %24, i8** %3)
  store i8 %25, i8* %6, align 1
  br label %39

26:                                               ; preds = %19
  %27 = load i8*, i8** %3, align 8
  %28 = call zeroext i8 @parse_int_type(i8* %27, i8** %3)
  store i8 %28, i8* %6, align 1
  br label %39

29:                                               ; preds = %19
  %30 = load i8*, i8** %3, align 8
  %31 = call zeroext i8 @parse_float_type(i8* %30, i8** %3)
  store i8 %31, i8* %6, align 1
  br label %39

32:                                               ; preds = %19
  %33 = load i8*, i8** %3, align 8
  %34 = call zeroext i8 @parse_bool_type(i8* %33, i8** %3)
  store i8 %34, i8* %6, align 1
  br label %39

35:                                               ; preds = %19
  %36 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %37 = load i8*, i8** %3, align 8
  %38 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %36, i8* getelementptr inbounds ([42 x i8], [42 x i8]* @.str.7, i64 0, i64 0), i8* %37)
  call void @exit(i32 127) #7
  unreachable

39:                                               ; preds = %32, %29, %26, %23
  br label %40

40:                                               ; preds = %39, %135, %137
  %41 = load i8*, i8** %3, align 8
  %42 = load i8, i8* %41, align 1
  %43 = sext i8 %42 to i32
  switch i32 %43, label %140 [
    i32 91, label %44
    i32 32, label %137
  ]

44:                                               ; preds = %40
  store i8 1, i8* %7, align 1
  %45 = load i8*, i8** %3, align 8
  %46 = getelementptr inbounds i8, i8* %45, i32 1
  store i8* %46, i8** %3, align 8
  br label %47

47:                                               ; preds = %52, %44
  %48 = load i8*, i8** %3, align 8
  %49 = load i8, i8* %48, align 1
  %50 = sext i8 %49 to i32
  %51 = icmp eq i32 %50, 32
  br i1 %51, label %52, label %55

52:                                               ; preds = %47
  %53 = load i8*, i8** %3, align 8
  %54 = getelementptr inbounds i8, i8* %53, i32 1
  store i8* %54, i8** %3, align 8
  br label %47

55:                                               ; preds = %47
  br label %56

56:                                               ; preds = %83, %55
  %57 = load i8*, i8** %3, align 8
  %58 = load i8, i8* %57, align 1
  %59 = sext i8 %58 to i32
  %60 = icmp eq i32 %59, 44
  br i1 %60, label %61, label %84

61:                                               ; preds = %56
  %62 = load i8, i8* %7, align 1
  %63 = zext i8 %62 to i32
  %64 = icmp sge i32 %63, 255
  br i1 %64, label %65, label %68

65:                                               ; preds = %61
  %66 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %67 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %66, i8* getelementptr inbounds ([45 x i8], [45 x i8]* @.str.8, i64 0, i64 0))
  call void @exit(i32 127) #7
  unreachable

68:                                               ; preds = %61
  %69 = load i8, i8* %7, align 1
  %70 = zext i8 %69 to i32
  %71 = add nsw i32 %70, 1
  %72 = trunc i32 %71 to i8
  store i8 %72, i8* %7, align 1
  %73 = load i8*, i8** %3, align 8
  %74 = getelementptr inbounds i8, i8* %73, i32 1
  store i8* %74, i8** %3, align 8
  br label %75

75:                                               ; preds = %80, %68
  %76 = load i8*, i8** %3, align 8
  %77 = load i8, i8* %76, align 1
  %78 = sext i8 %77 to i32
  %79 = icmp eq i32 %78, 32
  br i1 %79, label %80, label %83

80:                                               ; preds = %75
  %81 = load i8*, i8** %3, align 8
  %82 = getelementptr inbounds i8, i8* %81, i32 1
  store i8* %82, i8** %3, align 8
  br label %75

83:                                               ; preds = %75
  br label %56

84:                                               ; preds = %56
  br label %85

85:                                               ; preds = %90, %84
  %86 = load i8*, i8** %3, align 8
  %87 = load i8, i8* %86, align 1
  %88 = sext i8 %87 to i32
  %89 = icmp eq i32 %88, 32
  br i1 %89, label %90, label %93

90:                                               ; preds = %85
  %91 = load i8*, i8** %3, align 8
  %92 = getelementptr inbounds i8, i8* %91, i32 1
  store i8* %92, i8** %3, align 8
  br label %85

93:                                               ; preds = %85
  %94 = load i8*, i8** %3, align 8
  %95 = load i8, i8* %94, align 1
  %96 = sext i8 %95 to i32
  %97 = icmp ne i32 %96, 93
  br i1 %97, label %98, label %102

98:                                               ; preds = %93
  %99 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %100 = load i8*, i8** %5, align 8
  %101 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %99, i8* getelementptr inbounds ([42 x i8], [42 x i8]* @.str.7, i64 0, i64 0), i8* %100)
  call void @exit(i32 127) #7
  unreachable

102:                                              ; preds = %93
  %103 = load i8*, i8** %3, align 8
  %104 = getelementptr inbounds i8, i8* %103, i32 1
  store i8* %104, i8** %3, align 8
  %105 = load i8, i8* %7, align 1
  %106 = zext i8 %105 to i32
  %107 = icmp eq i32 %106, 1
  br i1 %107, label %108, label %119

108:                                              ; preds = %102
  %109 = call zeroext i8 @getmem(i64 2)
  store i8 %109, i8* %9, align 1
  store i8 %109, i8* %8, align 1
  %110 = load i8, i8* %9, align 1
  %111 = add i8 %110, 1
  store i8 %111, i8* %9, align 1
  %112 = zext i8 %110 to i64
  %113 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %112
  store i8 -2, i8* %113, align 1
  %114 = load i8, i8* %6, align 1
  %115 = load i8, i8* %9, align 1
  %116 = add i8 %115, 1
  store i8 %116, i8* %9, align 1
  %117 = zext i8 %115 to i64
  %118 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %117
  store i8 %114, i8* %118, align 1
  br label %135

119:                                              ; preds = %102
  %120 = call zeroext i8 @getmem(i64 3)
  store i8 %120, i8* %9, align 1
  store i8 %120, i8* %8, align 1
  %121 = load i8, i8* %9, align 1
  %122 = add i8 %121, 1
  store i8 %122, i8* %9, align 1
  %123 = zext i8 %121 to i64
  %124 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %123
  store i8 -1, i8* %124, align 1
  %125 = load i8, i8* %7, align 1
  %126 = load i8, i8* %9, align 1
  %127 = add i8 %126, 1
  store i8 %127, i8* %9, align 1
  %128 = zext i8 %126 to i64
  %129 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %128
  store i8 %125, i8* %129, align 1
  %130 = load i8, i8* %6, align 1
  %131 = load i8, i8* %9, align 1
  %132 = add i8 %131, 1
  store i8 %132, i8* %9, align 1
  %133 = zext i8 %131 to i64
  %134 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %133
  store i8 %130, i8* %134, align 1
  br label %135

135:                                              ; preds = %119, %108
  %136 = load i8, i8* %8, align 1
  store i8 %136, i8* %6, align 1
  br label %40

137:                                              ; preds = %40
  %138 = load i8*, i8** %3, align 8
  %139 = getelementptr inbounds i8, i8* %138, i32 1
  store i8* %139, i8** %3, align 8
  br label %40

140:                                              ; preds = %40
  %141 = load i8*, i8** %3, align 8
  %142 = load i8**, i8*** %4, align 8
  store i8* %141, i8** %142, align 8
  %143 = load i8, i8* %6, align 1
  ret i8 %143
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i64 @size_type(i8 zeroext %0) #0 {
  %2 = alloca i64, align 8
  %3 = alloca i8, align 1
  %4 = alloca i32, align 4
  %5 = alloca i32, align 4
  %6 = alloca i32, align 4
  %7 = alloca i32, align 4
  store i8 %0, i8* %3, align 1
  %8 = load i8, i8* %3, align 1
  %9 = zext i8 %8 to i64
  %10 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %9
  %11 = load i8, i8* %10, align 1
  %12 = zext i8 %11 to i32
  switch i32 %12, label %29 [
    i32 251, label %13
    i32 252, label %14
    i32 253, label %15
    i32 255, label %16
    i32 254, label %17
  ]

13:                                               ; preds = %1
  store i64 4, i64* %2, align 8
  br label %58

14:                                               ; preds = %1
  store i64 8, i64* %2, align 8
  br label %58

15:                                               ; preds = %1
  store i64 8, i64* %2, align 8
  br label %58

16:                                               ; preds = %1
  store i64 16, i64* %2, align 8
  br label %58

17:                                               ; preds = %1
  %18 = load i8, i8* %3, align 1
  %19 = zext i8 %18 to i32
  %20 = add nsw i32 %19, 1
  %21 = sext i32 %20 to i64
  %22 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %21
  %23 = load i8, i8* %22, align 1
  %24 = zext i8 %23 to i32
  store i32 %24, i32* %4, align 4
  %25 = load i32, i32* %4, align 4
  %26 = sext i32 %25 to i64
  %27 = mul i64 %26, 8
  %28 = add i64 %27, 8
  store i64 %28, i64* %2, align 8
  br label %58

29:                                               ; preds = %1
  %30 = load i8, i8* %3, align 1
  %31 = zext i8 %30 to i64
  %32 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %31
  %33 = load i8, i8* %32, align 1
  %34 = zext i8 %33 to i32
  store i32 %34, i32* %5, align 4
  store i32 0, i32* %6, align 4
  store i32 0, i32* %7, align 4
  br label %35

35:                                               ; preds = %52, %29
  %36 = load i32, i32* %7, align 4
  %37 = load i32, i32* %5, align 4
  %38 = icmp slt i32 %36, %37
  br i1 %38, label %39, label %55

39:                                               ; preds = %35
  %40 = load i8, i8* %3, align 1
  %41 = zext i8 %40 to i32
  %42 = load i32, i32* %7, align 4
  %43 = add nsw i32 %41, %42
  %44 = sext i32 %43 to i64
  %45 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %44
  %46 = load i8, i8* %45, align 1
  %47 = call i64 @size_type(i8 zeroext %46)
  %48 = load i32, i32* %6, align 4
  %49 = sext i32 %48 to i64
  %50 = add i64 %49, %47
  %51 = trunc i64 %50 to i32
  store i32 %51, i32* %6, align 4
  br label %52

52:                                               ; preds = %39
  %53 = load i32, i32* %7, align 4
  %54 = add nsw i32 %53, 1
  store i32 %54, i32* %7, align 4
  br label %35

55:                                               ; preds = %35
  %56 = load i32, i32* %6, align 4
  %57 = sext i32 %56 to i64
  store i64 %57, i64* %2, align 8
  br label %58

58:                                               ; preds = %55, %17, %16, %15, %14, %13
  %59 = load i64, i64* %2, align 8
  ret i64 %59
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @show_array(i8 zeroext %0, i32 %1, i64* %2) #0 {
  %4 = alloca i8, align 1
  %5 = alloca i32, align 4
  %6 = alloca i64*, align 8
  %7 = alloca i8*, align 8
  %8 = alloca i64, align 8
  %9 = alloca i32, align 4
  %10 = alloca i64, align 8
  %11 = alloca i64, align 8
  %12 = alloca i64, align 8
  %13 = alloca i64, align 8
  %14 = alloca i32, align 4
  store i8 %0, i8* %4, align 1
  store i32 %1, i32* %5, align 4
  store i64* %2, i64** %6, align 8
  %15 = load i64*, i64** %6, align 8
  %16 = load i32, i32* %5, align 4
  %17 = sext i32 %16 to i64
  %18 = getelementptr inbounds i64, i64* %15, i64 %17
  %19 = load i64, i64* %18, align 8
  %20 = inttoptr i64 %19 to i8*
  store i8* %20, i8** %7, align 8
  store i64 1, i64* %8, align 8
  store i32 0, i32* %9, align 4
  br label %21

21:                                               ; preds = %40, %3
  %22 = load i32, i32* %9, align 4
  %23 = load i32, i32* %5, align 4
  %24 = icmp slt i32 %22, %23
  br i1 %24, label %25, label %43

25:                                               ; preds = %21
  %26 = load i64*, i64** %6, align 8
  %27 = load i32, i32* %9, align 4
  %28 = sext i32 %27 to i64
  %29 = getelementptr inbounds i64, i64* %26, i64 %28
  %30 = load i64, i64* %29, align 8
  store i64 %30, i64* %10, align 8
  %31 = load i64, i64* %8, align 8
  %32 = load i64, i64* %10, align 8
  %33 = call { i64, i1 } @llvm.umul.with.overflow.i64(i64 %31, i64 %32)
  %34 = extractvalue { i64, i1 } %33, 1
  %35 = extractvalue { i64, i1 } %33, 0
  store i64 %35, i64* %8, align 8
  br i1 %34, label %36, label %39

36:                                               ; preds = %25
  %37 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %38 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %37, i8* getelementptr inbounds ([59 x i8], [59 x i8]* @.str.9, i64 0, i64 0))
  call void @exit(i32 127) #7
  unreachable

39:                                               ; preds = %25
  br label %40

40:                                               ; preds = %39
  %41 = load i32, i32* %9, align 4
  %42 = add nsw i32 %41, 1
  store i32 %42, i32* %9, align 4
  br label %21

43:                                               ; preds = %21
  %44 = load i8, i8* %4, align 1
  %45 = call i64 @size_type(i8 zeroext %44)
  store i64 %45, i64* %11, align 8
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.10, i64 0, i64 0))
  store i64 0, i64* %12, align 8
  br label %46

46:                                               ; preds = %106, %43
  %47 = load i64, i64* %12, align 8
  %48 = load i64, i64* %8, align 8
  %49 = icmp ult i64 %47, %48
  br i1 %49, label %50, label %109

50:                                               ; preds = %46
  %51 = load i8, i8* %4, align 1
  %52 = load i8*, i8** %7, align 8
  %53 = load i64, i64* %12, align 8
  %54 = load i64, i64* %11, align 8
  %55 = mul i64 %53, %54
  %56 = getelementptr i8, i8* %52, i64 %55
  call void @show_type(i8 zeroext %51, i8* %56)
  %57 = load i64, i64* %12, align 8
  %58 = add nsw i64 %57, 1
  store i64 %58, i64* %13, align 8
  store i32 0, i32* %14, align 4
  br label %59

59:                                               ; preds = %71, %50
  %60 = load i64, i64* %13, align 8
  %61 = load i64*, i64** %6, align 8
  %62 = load i32, i32* %5, align 4
  %63 = load i32, i32* %14, align 4
  %64 = sub nsw i32 %62, %63
  %65 = sub nsw i32 %64, 1
  %66 = sext i32 %65 to i64
  %67 = getelementptr inbounds i64, i64* %61, i64 %66
  %68 = load i64, i64* %67, align 8
  %69 = urem i64 %60, %68
  %70 = icmp eq i64 %69, 0
  br i1 %70, label %71, label %84

71:                                               ; preds = %59
  %72 = load i64*, i64** %6, align 8
  %73 = load i32, i32* %5, align 4
  %74 = load i32, i32* %14, align 4
  %75 = sub nsw i32 %73, %74
  %76 = sub nsw i32 %75, 1
  %77 = sext i32 %76 to i64
  %78 = getelementptr inbounds i64, i64* %72, i64 %77
  %79 = load i64, i64* %78, align 8
  %80 = load i64, i64* %13, align 8
  %81 = udiv i64 %80, %79
  store i64 %81, i64* %13, align 8
  %82 = load i32, i32* %14, align 4
  %83 = add nsw i32 %82, 1
  store i32 %83, i32* %14, align 4
  br label %59

84:                                               ; preds = %59
  %85 = load i64, i64* %12, align 8
  %86 = add nsw i64 %85, 1
  %87 = load i64, i64* %8, align 8
  %88 = icmp ult i64 %86, %87
  br i1 %88, label %89, label %105

89:                                               ; preds = %84
  %90 = load i32, i32* %14, align 4
  %91 = icmp eq i32 %90, 0
  br i1 %91, label %92, label %93

92:                                               ; preds = %89
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.11, i64 0, i64 0))
  br label %104

93:                                               ; preds = %89
  store i64 0, i64* %13, align 8
  br label %94

94:                                               ; preds = %100, %93
  %95 = load i64, i64* %13, align 8
  %96 = load i32, i32* %14, align 4
  %97 = sext i32 %96 to i64
  %98 = icmp slt i64 %95, %97
  br i1 %98, label %99, label %103

99:                                               ; preds = %94
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.12, i64 0, i64 0))
  br label %100

100:                                              ; preds = %99
  %101 = load i64, i64* %13, align 8
  %102 = add nsw i64 %101, 1
  store i64 %102, i64* %13, align 8
  br label %94

103:                                              ; preds = %94
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.13, i64 0, i64 0))
  br label %104

104:                                              ; preds = %103, %92
  br label %105

105:                                              ; preds = %104, %84
  br label %106

106:                                              ; preds = %105
  %107 = load i64, i64* %12, align 8
  %108 = add nsw i64 %107, 1
  store i64 %108, i64* %12, align 8
  br label %46

109:                                              ; preds = %46
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.14, i64 0, i64 0))
  ret void
}

; Function Attrs: nounwind readnone speculatable willreturn
declare { i64, i1 } @llvm.umul.with.overflow.i64(i64, i64) #4

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @show_type(i8 zeroext %0, i8* %1) #0 {
  %3 = alloca i8, align 1
  %4 = alloca i8*, align 8
  %5 = alloca i32, align 4
  %6 = alloca i32, align 4
  %7 = alloca i32, align 4
  store i8 %0, i8* %3, align 1
  store i8* %1, i8** %4, align 8
  %8 = load i8, i8* %3, align 1
  %9 = zext i8 %8 to i64
  %10 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %9
  %11 = load i8, i8* %10, align 1
  %12 = zext i8 %11 to i32
  switch i32 %12, label %54 [
    i32 251, label %13
    i32 252, label %21
    i32 253, label %25
    i32 254, label %29
    i32 255, label %38
  ]

13:                                               ; preds = %2
  %14 = load i8*, i8** %4, align 8
  %15 = bitcast i8* %14 to i32*
  %16 = load i32, i32* %15, align 4
  %17 = icmp ne i32 %16, 0
  br i1 %17, label %18, label %19

18:                                               ; preds = %13
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.15, i64 0, i64 0))
  br label %20

19:                                               ; preds = %13
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([6 x i8], [6 x i8]* @.str.16, i64 0, i64 0))
  br label %20

20:                                               ; preds = %19, %18
  br label %100

21:                                               ; preds = %2
  %22 = load i8*, i8** %4, align 8
  %23 = bitcast i8* %22 to i64*
  %24 = load i64, i64* %23, align 8
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.17, i64 0, i64 0), i64 %24)
  br label %100

25:                                               ; preds = %2
  %26 = load i8*, i8** %4, align 8
  %27 = bitcast i8* %26 to double*
  %28 = load double, double* %27, align 8
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.18, i64 0, i64 0), double %28)
  br label %100

29:                                               ; preds = %2
  %30 = load i8, i8* %3, align 1
  %31 = zext i8 %30 to i32
  %32 = add nsw i32 %31, 1
  %33 = sext i32 %32 to i64
  %34 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %33
  %35 = load i8, i8* %34, align 1
  %36 = load i8*, i8** %4, align 8
  %37 = bitcast i8* %36 to i64*
  call void @show_array(i8 zeroext %35, i32 1, i64* %37)
  br label %100

38:                                               ; preds = %2
  %39 = load i8, i8* %3, align 1
  %40 = zext i8 %39 to i32
  %41 = add nsw i32 %40, 2
  %42 = sext i32 %41 to i64
  %43 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %42
  %44 = load i8, i8* %43, align 1
  %45 = load i8, i8* %3, align 1
  %46 = zext i8 %45 to i32
  %47 = add nsw i32 %46, 1
  %48 = sext i32 %47 to i64
  %49 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %48
  %50 = load i8, i8* %49, align 1
  %51 = zext i8 %50 to i32
  %52 = load i8*, i8** %4, align 8
  %53 = bitcast i8* %52 to i64*
  call void @show_array(i8 zeroext %44, i32 %51, i64* %53)
  br label %100

54:                                               ; preds = %2
  %55 = load i8, i8* %3, align 1
  %56 = zext i8 %55 to i64
  %57 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %56
  %58 = load i8, i8* %57, align 1
  %59 = zext i8 %58 to i32
  store i32 %59, i32* %5, align 4
  store i32 0, i32* %6, align 4
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.19, i64 0, i64 0))
  store i32 0, i32* %7, align 4
  br label %60

60:                                               ; preds = %96, %54
  %61 = load i32, i32* %7, align 4
  %62 = load i32, i32* %5, align 4
  %63 = icmp slt i32 %61, %62
  br i1 %63, label %64, label %99

64:                                               ; preds = %60
  %65 = load i8, i8* %3, align 1
  %66 = zext i8 %65 to i32
  %67 = add nsw i32 %66, 1
  %68 = load i32, i32* %7, align 4
  %69 = add nsw i32 %67, %68
  %70 = sext i32 %69 to i64
  %71 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %70
  %72 = load i8, i8* %71, align 1
  %73 = load i8*, i8** %4, align 8
  %74 = load i32, i32* %6, align 4
  %75 = sext i32 %74 to i64
  %76 = getelementptr i8, i8* %73, i64 %75
  call void @show_type(i8 zeroext %72, i8* %76)
  %77 = load i8, i8* %3, align 1
  %78 = zext i8 %77 to i32
  %79 = add nsw i32 %78, 1
  %80 = load i32, i32* %7, align 4
  %81 = add nsw i32 %79, %80
  %82 = sext i32 %81 to i64
  %83 = getelementptr inbounds [256 x i8], [256 x i8]* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 1), i64 0, i64 %82
  %84 = load i8, i8* %83, align 1
  %85 = call i64 @size_type(i8 zeroext %84)
  %86 = load i32, i32* %6, align 4
  %87 = sext i32 %86 to i64
  %88 = add i64 %87, %85
  %89 = trunc i64 %88 to i32
  store i32 %89, i32* %6, align 4
  %90 = load i32, i32* %7, align 4
  %91 = add nsw i32 %90, 1
  %92 = load i32, i32* %5, align 4
  %93 = icmp ne i32 %91, %92
  br i1 %93, label %94, label %95

94:                                               ; preds = %64
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.11, i64 0, i64 0))
  br label %95

95:                                               ; preds = %94, %64
  br label %96

96:                                               ; preds = %95
  %97 = load i32, i32* %7, align 4
  %98 = add nsw i32 %97, 1
  store i32 %98, i32* %7, align 4
  br label %60

99:                                               ; preds = %60
  call void (i8*, ...) @tprintf(i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.20, i64 0, i64 0))
  br label %100

100:                                              ; preds = %99, %38, %29, %25, %21, %20
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @show(i8* %0, i8* %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8*, align 8
  %5 = alloca i8*, align 8
  %6 = alloca i8, align 1
  store i8* %0, i8** %3, align 8
  store i8* %1, i8** %4, align 8
  %7 = load i8*, i8** %3, align 8
  %8 = call i64 @strnlen(i8* %7, i64 256) #8
  %9 = icmp eq i64 %8, 256
  br i1 %9, label %10, label %14

10:                                               ; preds = %2
  %11 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %12 = load i8*, i8** %3, align 8
  %13 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %11, i8* getelementptr inbounds ([45 x i8], [45 x i8]* @.str.21, i64 0, i64 0), i8* %12)
  call void @exit(i32 127) #7
  unreachable

14:                                               ; preds = %2
  store i64 0, i64* getelementptr inbounds (%struct.anon, %struct.anon* @mem, i32 0, i32 0), align 8
  %15 = load i8*, i8** %3, align 8
  store i8* %15, i8** %5, align 8
  %16 = load i8*, i8** %3, align 8
  %17 = call zeroext i8 @parse_type(i8* %16, i8** %3)
  store i8 %17, i8* %6, align 1
  br label %18

18:                                               ; preds = %23, %14
  %19 = load i8*, i8** %3, align 8
  %20 = load i8, i8* %19, align 1
  %21 = sext i8 %20 to i32
  %22 = icmp eq i32 %21, 32
  br i1 %22, label %23, label %26

23:                                               ; preds = %18
  %24 = load i8*, i8** %3, align 8
  %25 = getelementptr inbounds i8, i8* %24, i32 1
  store i8* %25, i8** %3, align 8
  br label %18

26:                                               ; preds = %18
  %27 = load i8*, i8** %3, align 8
  %28 = load i8, i8* %27, align 1
  %29 = sext i8 %28 to i32
  %30 = icmp ne i32 %29, 0
  br i1 %30, label %31, label %35

31:                                               ; preds = %26
  %32 = load %struct._IO_FILE*, %struct._IO_FILE** @stderr, align 8
  %33 = load i8*, i8** %5, align 8
  %34 = call i32 (%struct._IO_FILE*, i8*, ...) @fprintf(%struct._IO_FILE* %32, i8* getelementptr inbounds ([46 x i8], [46 x i8]* @.str.22, i64 0, i64 0), i8* %33)
  call void @exit(i32 127) #7
  unreachable

35:                                               ; preds = %26
  %36 = load i8, i8* %6, align 1
  %37 = load i8*, i8** %4, align 8
  call void @show_type(i8 zeroext %36, i8* %37)
  ret i32 1
}

; Function Attrs: nounwind readonly
declare dso_local i64 @strnlen(i8*, i64) #5

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @_show(i8* %0, i8* %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8*, align 8
  store i8* %0, i8** %3, align 8
  store i8* %1, i8** %4, align 8
  %5 = load i8*, i8** %3, align 8
  %6 = load i8*, i8** %4, align 8
  %7 = call i32 @show(i8* %5, i8* %6)
  ret i32 %7
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @fail_assertion(i8* %0) #0 {
  %2 = alloca i8*, align 8
  store i8* %0, i8** %2, align 8
  %3 = load i8*, i8** %2, align 8
  %4 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([11 x i8], [11 x i8]* @.str.23, i64 0, i64 0), i8* %3)
  call void @exit(i32 1) #7
  unreachable
}

declare dso_local i32 @printf(i8*, ...) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @_fail_assertion(i8* %0) #0 {
  %2 = alloca i8*, align 8
  store i8* %0, i8** %2, align 8
  %3 = load i8*, i8** %2, align 8
  call void @fail_assertion(i8* %3)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @print(i8* %0) #0 {
  %2 = alloca i8*, align 8
  store i8* %0, i8** %2, align 8
  %3 = load i8*, i8** %2, align 8
  %4 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.24, i64 0, i64 0), i8* %3)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @_print(i8* %0) #0 {
  %2 = alloca i8*, align 8
  store i8* %0, i8** %2, align 8
  %3 = load i8*, i8** %2, align 8
  call void @print(i8* %3)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local double @get_time() #0 {
  %1 = alloca i64, align 8
  %2 = call i64 @clock() #3
  store i64 %2, i64* %1, align 8
  %3 = load i64, i64* %1, align 8
  %4 = sitofp i64 %3 to double
  %5 = fdiv double %4, 1.000000e+06
  ret double %5
}

; Function Attrs: nounwind
declare dso_local i64 @clock() #6

; Function Attrs: noinline nounwind optnone uwtable
define dso_local double @_get_time() #0 {
  %1 = call double @get_time()
  ret double %1
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i64 @sub_ints(i64 %0, i64 %1) #0 {
  %3 = alloca i64, align 8
  %4 = alloca i64, align 8
  store i64 %0, i64* %3, align 8
  store i64 %1, i64* %4, align 8
  %5 = load i64, i64* %3, align 8
  %6 = load i64, i64* %4, align 8
  %7 = sub nsw i64 %5, %6
  ret i64 %7
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i64 @_sub_ints(i64 %0, i64 %1) #0 {
  %3 = alloca i64, align 8
  %4 = alloca i64, align 8
  store i64 %0, i64* %3, align 8
  store i64 %1, i64* %4, align 8
  %5 = load i64, i64* %3, align 8
  %6 = load i64, i64* %4, align 8
  %7 = call i64 @sub_ints(i64 %5, i64 %6)
  ret i64 %7
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local double @sub_floats(double %0, double %1) #0 {
  %3 = alloca double, align 8
  %4 = alloca double, align 8
  store double %0, double* %3, align 8
  store double %1, double* %4, align 8
  %5 = load double, double* %3, align 8
  %6 = load double, double* %4, align 8
  %7 = fsub double %5, %6
  ret double %7
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local double @_sub_floats(double %0, double %1) #0 {
  %3 = alloca double, align 8
  %4 = alloca double, align 8
  store double %0, double* %3, align 8
  store double %1, double* %4, align 8
  %5 = load double, double* %3, align 8
  %6 = load double, double* %4, align 8
  %7 = call double @sub_floats(double %5, double %6)
  ret double %7
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @has_size(%struct.pict* byval(%struct.pict) align 8 %0, i64 %1, i64 %2) #0 {
  %4 = alloca i64, align 8
  %5 = alloca i64, align 8
  store i64 %1, i64* %4, align 8
  store i64 %2, i64* %5, align 8
  %6 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 0
  %7 = load i64, i64* %6, align 8
  %8 = load i64, i64* %4, align 8
  %9 = icmp eq i64 %7, %8
  br i1 %9, label %10, label %15

10:                                               ; preds = %3
  %11 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %12 = load i64, i64* %11, align 8
  %13 = load i64, i64* %5, align 8
  %14 = icmp eq i64 %12, %13
  br label %15

15:                                               ; preds = %10, %3
  %16 = phi i1 [ false, %3 ], [ %14, %10 ]
  %17 = zext i1 %16 to i32
  ret i32 %17
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @_has_size(%struct.pict* byval(%struct.pict) align 8 %0, i64 %1, i64 %2) #0 {
  %4 = alloca i64, align 8
  %5 = alloca i64, align 8
  store i64 %1, i64* %4, align 8
  store i64 %2, i64* %5, align 8
  %6 = load i64, i64* %4, align 8
  %7 = load i64, i64* %5, align 8
  %8 = call i32 @has_size(%struct.pict* byval(%struct.pict) align 8 %0, i64 %6, i64 %7)
  ret i32 %8
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @sepia(%struct.pict* noalias sret %0, %struct.pict* byval(%struct.pict) align 8 %1) #0 {
  %3 = alloca i64, align 8
  %4 = alloca i64, align 8
  %5 = alloca double, align 8
  %6 = alloca double, align 8
  %7 = alloca double, align 8
  %8 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %9 = load i64, i64* %8, align 8
  %10 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 0
  store i64 %9, i64* %10, align 8
  %11 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %12 = load i64, i64* %11, align 8
  %13 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  store i64 %12, i64* %13, align 8
  %14 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %15 = load i64, i64* %14, align 8
  %16 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %17 = load i64, i64* %16, align 8
  %18 = mul nsw i64 %15, %17
  %19 = mul nsw i64 %18, 4
  %20 = mul i64 %19, 8
  %21 = call noalias i8* @malloc(i64 %20) #3
  %22 = bitcast i8* %21 to double*
  %23 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  store double* %22, double** %23, align 8
  %24 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %25 = load double*, double** %24, align 8
  %26 = icmp ne double* %25, null
  br i1 %26, label %28, label %27

27:                                               ; preds = %2
  call void @fail_assertion(i8* getelementptr inbounds ([15 x i8], [15 x i8]* @.str.25, i64 0, i64 0))
  br label %28

28:                                               ; preds = %27, %2
  store i64 0, i64* %3, align 8
  br label %29

29:                                               ; preds = %149, %28
  %30 = load i64, i64* %3, align 8
  %31 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %32 = load i64, i64* %31, align 8
  %33 = icmp slt i64 %30, %32
  br i1 %33, label %34, label %152

34:                                               ; preds = %29
  store i64 0, i64* %4, align 8
  br label %35

35:                                               ; preds = %145, %34
  %36 = load i64, i64* %4, align 8
  %37 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %38 = load i64, i64* %37, align 8
  %39 = icmp slt i64 %36, %38
  br i1 %39, label %40, label %148

40:                                               ; preds = %35
  %41 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %42 = load double*, double** %41, align 8
  %43 = load i64, i64* %3, align 8
  %44 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %45 = load i64, i64* %44, align 8
  %46 = mul nsw i64 %43, %45
  %47 = load i64, i64* %4, align 8
  %48 = add nsw i64 %46, %47
  %49 = mul nsw i64 4, %48
  %50 = add nsw i64 %49, 0
  %51 = getelementptr inbounds double, double* %42, i64 %50
  %52 = load double, double* %51, align 8
  store double %52, double* %5, align 8
  %53 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %54 = load double*, double** %53, align 8
  %55 = load i64, i64* %3, align 8
  %56 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %57 = load i64, i64* %56, align 8
  %58 = mul nsw i64 %55, %57
  %59 = load i64, i64* %4, align 8
  %60 = add nsw i64 %58, %59
  %61 = mul nsw i64 4, %60
  %62 = add nsw i64 %61, 1
  %63 = getelementptr inbounds double, double* %54, i64 %62
  %64 = load double, double* %63, align 8
  store double %64, double* %6, align 8
  %65 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %66 = load double*, double** %65, align 8
  %67 = load i64, i64* %3, align 8
  %68 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %69 = load i64, i64* %68, align 8
  %70 = mul nsw i64 %67, %69
  %71 = load i64, i64* %4, align 8
  %72 = add nsw i64 %70, %71
  %73 = mul nsw i64 4, %72
  %74 = add nsw i64 %73, 2
  %75 = getelementptr inbounds double, double* %66, i64 %74
  %76 = load double, double* %75, align 8
  store double %76, double* %7, align 8
  %77 = load double, double* %5, align 8
  %78 = fmul double 3.930000e-01, %77
  %79 = load double, double* %6, align 8
  %80 = fmul double 7.690000e-01, %79
  %81 = fadd double %78, %80
  %82 = load double, double* %7, align 8
  %83 = fmul double 1.890000e-01, %82
  %84 = fadd double %81, %83
  %85 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %86 = load double*, double** %85, align 8
  %87 = load i64, i64* %3, align 8
  %88 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %89 = load i64, i64* %88, align 8
  %90 = mul nsw i64 %87, %89
  %91 = load i64, i64* %4, align 8
  %92 = add nsw i64 %90, %91
  %93 = mul nsw i64 4, %92
  %94 = add nsw i64 %93, 0
  %95 = getelementptr inbounds double, double* %86, i64 %94
  store double %84, double* %95, align 8
  %96 = load double, double* %5, align 8
  %97 = fmul double 3.490000e-01, %96
  %98 = load double, double* %6, align 8
  %99 = fmul double 6.860000e-01, %98
  %100 = fadd double %97, %99
  %101 = load double, double* %7, align 8
  %102 = fmul double 1.680000e-01, %101
  %103 = fadd double %100, %102
  %104 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %105 = load double*, double** %104, align 8
  %106 = load i64, i64* %3, align 8
  %107 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %108 = load i64, i64* %107, align 8
  %109 = mul nsw i64 %106, %108
  %110 = load i64, i64* %4, align 8
  %111 = add nsw i64 %109, %110
  %112 = mul nsw i64 4, %111
  %113 = add nsw i64 %112, 1
  %114 = getelementptr inbounds double, double* %105, i64 %113
  store double %103, double* %114, align 8
  %115 = load double, double* %5, align 8
  %116 = fmul double 2.720000e-01, %115
  %117 = load double, double* %6, align 8
  %118 = fmul double 5.340000e-01, %117
  %119 = fadd double %116, %118
  %120 = load double, double* %7, align 8
  %121 = fmul double 1.310000e-01, %120
  %122 = fadd double %119, %121
  %123 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %124 = load double*, double** %123, align 8
  %125 = load i64, i64* %3, align 8
  %126 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %127 = load i64, i64* %126, align 8
  %128 = mul nsw i64 %125, %127
  %129 = load i64, i64* %4, align 8
  %130 = add nsw i64 %128, %129
  %131 = mul nsw i64 4, %130
  %132 = add nsw i64 %131, 2
  %133 = getelementptr inbounds double, double* %124, i64 %132
  store double %122, double* %133, align 8
  %134 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %135 = load double*, double** %134, align 8
  %136 = load i64, i64* %3, align 8
  %137 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %138 = load i64, i64* %137, align 8
  %139 = mul nsw i64 %136, %138
  %140 = load i64, i64* %4, align 8
  %141 = add nsw i64 %139, %140
  %142 = mul nsw i64 4, %141
  %143 = add nsw i64 %142, 3
  %144 = getelementptr inbounds double, double* %135, i64 %143
  store double 1.000000e+00, double* %144, align 8
  br label %145

145:                                              ; preds = %40
  %146 = load i64, i64* %4, align 8
  %147 = add nsw i64 %146, 1
  store i64 %147, i64* %4, align 8
  br label %35

148:                                              ; preds = %35
  br label %149

149:                                              ; preds = %148
  %150 = load i64, i64* %3, align 8
  %151 = add nsw i64 %150, 1
  store i64 %151, i64* %3, align 8
  br label %29

152:                                              ; preds = %29
  ret void
}

; Function Attrs: nounwind
declare dso_local noalias i8* @malloc(i64) #6

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @_sepia(%struct.pict* noalias sret %0, %struct.pict* byval(%struct.pict) align 8 %1) #0 {
  call void @sepia(%struct.pict* sret %0, %struct.pict* byval(%struct.pict) align 8 %1)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @blur(%struct.pict* noalias sret %0, %struct.pict* byval(%struct.pict) align 8 %1, double %2) #0 {
  %4 = alloca double, align 8
  %5 = alloca i64, align 8
  %6 = alloca i64, align 8
  %7 = alloca double*, align 8
  %8 = alloca double, align 8
  %9 = alloca double, align 8
  %10 = alloca i64, align 8
  %11 = alloca i64, align 8
  %12 = alloca double, align 8
  %13 = alloca i64, align 8
  %14 = alloca i64, align 8
  %15 = alloca double, align 8
  %16 = alloca double, align 8
  %17 = alloca double, align 8
  %18 = alloca double, align 8
  %19 = alloca i64, align 8
  %20 = alloca i64, align 8
  %21 = alloca double, align 8
  store double %2, double* %4, align 8
  %22 = load double, double* %4, align 8
  %23 = fcmp ole double %22, 0.000000e+00
  br i1 %23, label %24, label %25

24:                                               ; preds = %3
  call void @fail_assertion(i8* getelementptr inbounds ([30 x i8], [30 x i8]* @.str.26, i64 0, i64 0))
  br label %25

25:                                               ; preds = %24, %3
  %26 = load double, double* %4, align 8
  %27 = fmul double 3.000000e+00, %26
  %28 = fadd double %27, 5.000000e-01
  %29 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %30 = load i64, i64* %29, align 8
  %31 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %32 = load i64, i64* %31, align 8
  %33 = icmp slt i64 %30, %32
  br i1 %33, label %34, label %37

34:                                               ; preds = %25
  %35 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %36 = load i64, i64* %35, align 8
  br label %40

37:                                               ; preds = %25
  %38 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %39 = load i64, i64* %38, align 8
  br label %40

40:                                               ; preds = %37, %34
  %41 = phi i64 [ %36, %34 ], [ %39, %37 ]
  %42 = sitofp i64 %41 to double
  %43 = call double @llvm.minnum.f64(double %28, double %42)
  %44 = fptosi double %43 to i32
  %45 = sext i32 %44 to i64
  store i64 %45, i64* %5, align 8
  %46 = load i64, i64* %5, align 8
  %47 = mul nsw i64 2, %46
  %48 = add nsw i64 %47, 1
  store i64 %48, i64* %6, align 8
  %49 = load i64, i64* %6, align 8
  %50 = load i64, i64* %6, align 8
  %51 = mul nsw i64 %49, %50
  %52 = mul i64 %51, 8
  %53 = call noalias i8* @malloc(i64 %52) #3
  %54 = bitcast i8* %53 to double*
  store double* %54, double** %7, align 8
  %55 = load double*, double** %7, align 8
  %56 = icmp ne double* %55, null
  br i1 %56, label %58, label %57

57:                                               ; preds = %40
  call void @fail_assertion(i8* getelementptr inbounds ([15 x i8], [15 x i8]* @.str.25, i64 0, i64 0))
  br label %58

58:                                               ; preds = %57, %40
  %59 = load double, double* %4, align 8
  %60 = fmul double 0x401921FB54442D18, %59
  %61 = load double, double* %4, align 8
  %62 = fmul double %60, %61
  %63 = fdiv double 1.000000e+00, %62
  store double %63, double* %8, align 8
  %64 = load double, double* %4, align 8
  %65 = fmul double 2.000000e+00, %64
  %66 = load double, double* %4, align 8
  %67 = fmul double %65, %66
  %68 = fdiv double -1.000000e+00, %67
  store double %68, double* %9, align 8
  %69 = load i64, i64* %5, align 8
  %70 = sub nsw i64 0, %69
  store i64 %70, i64* %10, align 8
  br label %71

71:                                               ; preds = %112, %58
  %72 = load i64, i64* %10, align 8
  %73 = load i64, i64* %5, align 8
  %74 = icmp sle i64 %72, %73
  br i1 %74, label %75, label %115

75:                                               ; preds = %71
  %76 = load i64, i64* %5, align 8
  %77 = sub nsw i64 0, %76
  store i64 %77, i64* %11, align 8
  br label %78

78:                                               ; preds = %108, %75
  %79 = load i64, i64* %11, align 8
  %80 = load i64, i64* %5, align 8
  %81 = icmp sle i64 %79, %80
  br i1 %81, label %82, label %111

82:                                               ; preds = %78
  %83 = load i64, i64* %10, align 8
  %84 = load i64, i64* %10, align 8
  %85 = mul nsw i64 %83, %84
  %86 = load i64, i64* %11, align 8
  %87 = load i64, i64* %11, align 8
  %88 = mul nsw i64 %86, %87
  %89 = add nsw i64 %85, %88
  %90 = sitofp i64 %89 to double
  store double %90, double* %12, align 8
  %91 = load double, double* %8, align 8
  %92 = load double, double* %12, align 8
  %93 = load double, double* %9, align 8
  %94 = fmul double %92, %93
  %95 = call double @exp(double %94) #3
  %96 = fmul double %91, %95
  %97 = load double*, double** %7, align 8
  %98 = load i64, i64* %10, align 8
  %99 = load i64, i64* %5, align 8
  %100 = add nsw i64 %98, %99
  %101 = load i64, i64* %6, align 8
  %102 = mul nsw i64 %100, %101
  %103 = load i64, i64* %11, align 8
  %104 = add nsw i64 %102, %103
  %105 = load i64, i64* %5, align 8
  %106 = add nsw i64 %104, %105
  %107 = getelementptr inbounds double, double* %97, i64 %106
  store double %96, double* %107, align 8
  br label %108

108:                                              ; preds = %82
  %109 = load i64, i64* %11, align 8
  %110 = add nsw i64 %109, 1
  store i64 %110, i64* %11, align 8
  br label %78

111:                                              ; preds = %78
  br label %112

112:                                              ; preds = %111
  %113 = load i64, i64* %10, align 8
  %114 = add nsw i64 %113, 1
  store i64 %114, i64* %10, align 8
  br label %71

115:                                              ; preds = %71
  %116 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %117 = load i64, i64* %116, align 8
  %118 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 0
  store i64 %117, i64* %118, align 8
  %119 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %120 = load i64, i64* %119, align 8
  %121 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  store i64 %120, i64* %121, align 8
  %122 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %123 = load i64, i64* %122, align 8
  %124 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %125 = load i64, i64* %124, align 8
  %126 = mul nsw i64 %123, %125
  %127 = mul nsw i64 %126, 4
  %128 = mul i64 %127, 8
  %129 = call noalias i8* @malloc(i64 %128) #3
  %130 = bitcast i8* %129 to double*
  %131 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  store double* %130, double** %131, align 8
  %132 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %133 = load double*, double** %132, align 8
  %134 = icmp ne double* %133, null
  br i1 %134, label %136, label %135

135:                                              ; preds = %115
  call void @fail_assertion(i8* getelementptr inbounds ([15 x i8], [15 x i8]* @.str.25, i64 0, i64 0))
  br label %136

136:                                              ; preds = %135, %115
  store i64 0, i64* %13, align 8
  br label %137

137:                                              ; preds = %321, %136
  %138 = load i64, i64* %13, align 8
  %139 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %140 = load i64, i64* %139, align 8
  %141 = icmp slt i64 %138, %140
  br i1 %141, label %142, label %324

142:                                              ; preds = %137
  store i64 0, i64* %14, align 8
  br label %143

143:                                              ; preds = %317, %142
  %144 = load i64, i64* %14, align 8
  %145 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %146 = load i64, i64* %145, align 8
  %147 = icmp slt i64 %144, %146
  br i1 %147, label %148, label %320

148:                                              ; preds = %143
  store double 0.000000e+00, double* %15, align 8
  store double 0.000000e+00, double* %16, align 8
  store double 0.000000e+00, double* %17, align 8
  store double 0.000000e+00, double* %18, align 8
  %149 = load i64, i64* %13, align 8
  %150 = load i64, i64* %5, align 8
  %151 = sub nsw i64 %149, %150
  store i64 %151, i64* %19, align 8
  br label %152

152:                                              ; preds = %260, %148
  %153 = load i64, i64* %19, align 8
  %154 = load i64, i64* %13, align 8
  %155 = load i64, i64* %5, align 8
  %156 = add nsw i64 %154, %155
  %157 = icmp sle i64 %153, %156
  br i1 %157, label %158, label %263

158:                                              ; preds = %152
  %159 = load i64, i64* %19, align 8
  %160 = icmp slt i64 %159, 0
  br i1 %160, label %161, label %162

161:                                              ; preds = %158
  br label %260

162:                                              ; preds = %158
  %163 = load i64, i64* %19, align 8
  %164 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %165 = load i64, i64* %164, align 8
  %166 = icmp sge i64 %163, %165
  br i1 %166, label %167, label %168

167:                                              ; preds = %162
  br label %260

168:                                              ; preds = %162
  %169 = load i64, i64* %14, align 8
  %170 = load i64, i64* %5, align 8
  %171 = sub nsw i64 %169, %170
  store i64 %171, i64* %20, align 8
  br label %172

172:                                              ; preds = %256, %168
  %173 = load i64, i64* %20, align 8
  %174 = load i64, i64* %14, align 8
  %175 = load i64, i64* %5, align 8
  %176 = add nsw i64 %174, %175
  %177 = icmp sle i64 %173, %176
  br i1 %177, label %178, label %259

178:                                              ; preds = %172
  %179 = load i64, i64* %20, align 8
  %180 = icmp slt i64 %179, 0
  br i1 %180, label %181, label %182

181:                                              ; preds = %178
  br label %256

182:                                              ; preds = %178
  %183 = load i64, i64* %20, align 8
  %184 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %185 = load i64, i64* %184, align 8
  %186 = icmp sge i64 %183, %185
  br i1 %186, label %187, label %188

187:                                              ; preds = %182
  br label %256

188:                                              ; preds = %182
  %189 = load double*, double** %7, align 8
  %190 = load i64, i64* %6, align 8
  %191 = load i64, i64* %19, align 8
  %192 = load i64, i64* %13, align 8
  %193 = sub nsw i64 %191, %192
  %194 = load i64, i64* %5, align 8
  %195 = add nsw i64 %193, %194
  %196 = mul nsw i64 %190, %195
  %197 = load i64, i64* %20, align 8
  %198 = load i64, i64* %14, align 8
  %199 = sub nsw i64 %197, %198
  %200 = load i64, i64* %5, align 8
  %201 = add nsw i64 %199, %200
  %202 = add nsw i64 %196, %201
  %203 = getelementptr inbounds double, double* %189, i64 %202
  %204 = load double, double* %203, align 8
  store double %204, double* %21, align 8
  %205 = load double, double* %21, align 8
  %206 = load double, double* %18, align 8
  %207 = fadd double %206, %205
  store double %207, double* %18, align 8
  %208 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %209 = load double*, double** %208, align 8
  %210 = load i64, i64* %19, align 8
  %211 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %212 = load i64, i64* %211, align 8
  %213 = mul nsw i64 %210, %212
  %214 = load i64, i64* %20, align 8
  %215 = add nsw i64 %213, %214
  %216 = mul nsw i64 4, %215
  %217 = add nsw i64 %216, 0
  %218 = getelementptr inbounds double, double* %209, i64 %217
  %219 = load double, double* %218, align 8
  %220 = load double, double* %21, align 8
  %221 = fmul double %219, %220
  %222 = load double, double* %15, align 8
  %223 = fadd double %222, %221
  store double %223, double* %15, align 8
  %224 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %225 = load double*, double** %224, align 8
  %226 = load i64, i64* %19, align 8
  %227 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %228 = load i64, i64* %227, align 8
  %229 = mul nsw i64 %226, %228
  %230 = load i64, i64* %20, align 8
  %231 = add nsw i64 %229, %230
  %232 = mul nsw i64 4, %231
  %233 = add nsw i64 %232, 1
  %234 = getelementptr inbounds double, double* %225, i64 %233
  %235 = load double, double* %234, align 8
  %236 = load double, double* %21, align 8
  %237 = fmul double %235, %236
  %238 = load double, double* %16, align 8
  %239 = fadd double %238, %237
  store double %239, double* %16, align 8
  %240 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %241 = load double*, double** %240, align 8
  %242 = load i64, i64* %19, align 8
  %243 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %244 = load i64, i64* %243, align 8
  %245 = mul nsw i64 %242, %244
  %246 = load i64, i64* %20, align 8
  %247 = add nsw i64 %245, %246
  %248 = mul nsw i64 4, %247
  %249 = add nsw i64 %248, 2
  %250 = getelementptr inbounds double, double* %241, i64 %249
  %251 = load double, double* %250, align 8
  %252 = load double, double* %21, align 8
  %253 = fmul double %251, %252
  %254 = load double, double* %17, align 8
  %255 = fadd double %254, %253
  store double %255, double* %17, align 8
  br label %256

256:                                              ; preds = %188, %187, %181
  %257 = load i64, i64* %20, align 8
  %258 = add nsw i64 %257, 1
  store i64 %258, i64* %20, align 8
  br label %172

259:                                              ; preds = %172
  br label %260

260:                                              ; preds = %259, %167, %161
  %261 = load i64, i64* %19, align 8
  %262 = add nsw i64 %261, 1
  store i64 %262, i64* %19, align 8
  br label %152

263:                                              ; preds = %152
  %264 = load double, double* %15, align 8
  %265 = load double, double* %18, align 8
  %266 = fdiv double %264, %265
  %267 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %268 = load double*, double** %267, align 8
  %269 = load i64, i64* %13, align 8
  %270 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %271 = load i64, i64* %270, align 8
  %272 = mul nsw i64 %269, %271
  %273 = load i64, i64* %14, align 8
  %274 = add nsw i64 %272, %273
  %275 = mul nsw i64 4, %274
  %276 = add nsw i64 %275, 0
  %277 = getelementptr inbounds double, double* %268, i64 %276
  store double %266, double* %277, align 8
  %278 = load double, double* %16, align 8
  %279 = load double, double* %18, align 8
  %280 = fdiv double %278, %279
  %281 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %282 = load double*, double** %281, align 8
  %283 = load i64, i64* %13, align 8
  %284 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %285 = load i64, i64* %284, align 8
  %286 = mul nsw i64 %283, %285
  %287 = load i64, i64* %14, align 8
  %288 = add nsw i64 %286, %287
  %289 = mul nsw i64 4, %288
  %290 = add nsw i64 %289, 1
  %291 = getelementptr inbounds double, double* %282, i64 %290
  store double %280, double* %291, align 8
  %292 = load double, double* %17, align 8
  %293 = load double, double* %18, align 8
  %294 = fdiv double %292, %293
  %295 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %296 = load double*, double** %295, align 8
  %297 = load i64, i64* %13, align 8
  %298 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %299 = load i64, i64* %298, align 8
  %300 = mul nsw i64 %297, %299
  %301 = load i64, i64* %14, align 8
  %302 = add nsw i64 %300, %301
  %303 = mul nsw i64 4, %302
  %304 = add nsw i64 %303, 2
  %305 = getelementptr inbounds double, double* %296, i64 %304
  store double %294, double* %305, align 8
  %306 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %307 = load double*, double** %306, align 8
  %308 = load i64, i64* %13, align 8
  %309 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %310 = load i64, i64* %309, align 8
  %311 = mul nsw i64 %308, %310
  %312 = load i64, i64* %14, align 8
  %313 = add nsw i64 %311, %312
  %314 = mul nsw i64 4, %313
  %315 = add nsw i64 %314, 3
  %316 = getelementptr inbounds double, double* %307, i64 %315
  store double 1.000000e+00, double* %316, align 8
  br label %317

317:                                              ; preds = %263
  %318 = load i64, i64* %14, align 8
  %319 = add nsw i64 %318, 1
  store i64 %319, i64* %14, align 8
  br label %143

320:                                              ; preds = %143
  br label %321

321:                                              ; preds = %320
  %322 = load i64, i64* %13, align 8
  %323 = add nsw i64 %322, 1
  store i64 %323, i64* %13, align 8
  br label %137

324:                                              ; preds = %137
  ret void
}

; Function Attrs: nounwind readnone speculatable willreturn
declare double @llvm.minnum.f64(double, double) #4

; Function Attrs: nounwind
declare dso_local double @exp(double) #6

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @_blur(%struct.pict* noalias sret %0, %struct.pict* byval(%struct.pict) align 8 %1, double %2) #0 {
  %4 = alloca double, align 8
  store double %2, double* %4, align 8
  %5 = load double, double* %4, align 8
  call void @blur(%struct.pict* sret %0, %struct.pict* byval(%struct.pict) align 8 %1, double %5)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @resize(%struct.pict* noalias sret %0, %struct.pict* byval(%struct.pict) align 8 %1, i32 %2, i32 %3) #0 {
  %5 = alloca i32, align 4
  %6 = alloca i32, align 4
  %7 = alloca i64, align 8
  %8 = alloca i64, align 8
  %9 = alloca double, align 8
  %10 = alloca double, align 8
  %11 = alloca i64, align 8
  %12 = alloca i64, align 8
  %13 = alloca double, align 8
  %14 = alloca double, align 8
  %15 = alloca double, align 8
  %16 = alloca double, align 8
  %17 = alloca double, align 8
  %18 = alloca double, align 8
  store i32 %2, i32* %5, align 4
  store i32 %3, i32* %6, align 4
  %19 = load i32, i32* %5, align 4
  %20 = icmp sle i32 %19, 0
  br i1 %20, label %24, label %21

21:                                               ; preds = %4
  %22 = load i32, i32* %6, align 4
  %23 = icmp sle i32 %22, 0
  br i1 %23, label %24, label %25

24:                                               ; preds = %21, %4
  call void @fail_assertion(i8* getelementptr inbounds ([32 x i8], [32 x i8]* @.str.27, i64 0, i64 0))
  br label %25

25:                                               ; preds = %24, %21
  %26 = load i32, i32* %5, align 4
  %27 = sext i32 %26 to i64
  %28 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 0
  store i64 %27, i64* %28, align 8
  %29 = load i32, i32* %6, align 4
  %30 = sext i32 %29 to i64
  %31 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  store i64 %30, i64* %31, align 8
  %32 = load i32, i32* %5, align 4
  %33 = load i32, i32* %6, align 4
  %34 = mul nsw i32 %32, %33
  %35 = mul nsw i32 %34, 4
  %36 = sext i32 %35 to i64
  %37 = mul i64 %36, 8
  %38 = call noalias i8* @malloc(i64 %37) #3
  %39 = bitcast i8* %38 to double*
  %40 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  store double* %39, double** %40, align 8
  %41 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %42 = load double*, double** %41, align 8
  %43 = icmp ne double* %42, null
  br i1 %43, label %45, label %44

44:                                               ; preds = %25
  call void @fail_assertion(i8* getelementptr inbounds ([15 x i8], [15 x i8]* @.str.25, i64 0, i64 0))
  br label %45

45:                                               ; preds = %44, %25
  store i64 0, i64* %7, align 8
  br label %46

46:                                               ; preds = %441, %45
  %47 = load i64, i64* %7, align 8
  %48 = load i32, i32* %5, align 4
  %49 = sext i32 %48 to i64
  %50 = icmp slt i64 %47, %49
  br i1 %50, label %51, label %444

51:                                               ; preds = %46
  store i64 0, i64* %8, align 8
  br label %52

52:                                               ; preds = %437, %51
  %53 = load i64, i64* %8, align 8
  %54 = load i32, i32* %6, align 4
  %55 = sext i32 %54 to i64
  %56 = icmp slt i64 %53, %55
  br i1 %56, label %57, label %440

57:                                               ; preds = %52
  %58 = load i64, i64* %7, align 8
  %59 = sitofp i64 %58 to double
  %60 = fadd double %59, 5.000000e-01
  %61 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %62 = load i64, i64* %61, align 8
  %63 = sitofp i64 %62 to double
  %64 = fmul double %60, %63
  %65 = load i32, i32* %5, align 4
  %66 = sitofp i32 %65 to double
  %67 = fdiv double %64, %66
  %68 = fsub double %67, 5.000000e-01
  store double %68, double* %9, align 8
  %69 = load i64, i64* %8, align 8
  %70 = sitofp i64 %69 to double
  %71 = fadd double %70, 5.000000e-01
  %72 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 0
  %73 = load i64, i64* %72, align 8
  %74 = sitofp i64 %73 to double
  %75 = fmul double %71, %74
  %76 = load i32, i32* %5, align 4
  %77 = sitofp i32 %76 to double
  %78 = fdiv double %75, %77
  %79 = fsub double %78, 5.000000e-01
  store double %79, double* %10, align 8
  %80 = load double, double* %9, align 8
  %81 = fptosi double %80 to i64
  store i64 %81, i64* %11, align 8
  %82 = load double, double* %10, align 8
  %83 = fptosi double %82 to i64
  store i64 %83, i64* %12, align 8
  %84 = load double, double* %9, align 8
  %85 = load i64, i64* %11, align 8
  %86 = sitofp i64 %85 to double
  %87 = fsub double %84, %86
  store double %87, double* %13, align 8
  %88 = load double, double* %10, align 8
  %89 = load i64, i64* %12, align 8
  %90 = sitofp i64 %89 to double
  %91 = fsub double %88, %90
  store double %91, double* %14, align 8
  %92 = load double, double* %13, align 8
  %93 = fsub double 1.000000e+00, %92
  %94 = load double, double* %14, align 8
  %95 = fsub double 1.000000e+00, %94
  %96 = fmul double %93, %95
  %97 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %98 = load double*, double** %97, align 8
  %99 = load i64, i64* %11, align 8
  %100 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %101 = load i64, i64* %100, align 8
  %102 = mul nsw i64 %99, %101
  %103 = load i64, i64* %12, align 8
  %104 = add nsw i64 %102, %103
  %105 = mul nsw i64 4, %104
  %106 = add nsw i64 %105, 0
  %107 = getelementptr inbounds double, double* %98, i64 %106
  %108 = load double, double* %107, align 8
  %109 = fmul double %96, %108
  %110 = load double, double* %13, align 8
  %111 = fcmp ogt double %110, 0.000000e+00
  br i1 %111, label %112, label %132

112:                                              ; preds = %57
  %113 = load double, double* %13, align 8
  %114 = fadd double 0.000000e+00, %113
  %115 = load double, double* %14, align 8
  %116 = fsub double 1.000000e+00, %115
  %117 = fmul double %114, %116
  %118 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %119 = load double*, double** %118, align 8
  %120 = load i64, i64* %11, align 8
  %121 = add nsw i64 %120, 1
  %122 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %123 = load i64, i64* %122, align 8
  %124 = mul nsw i64 %121, %123
  %125 = load i64, i64* %12, align 8
  %126 = add nsw i64 %124, %125
  %127 = mul nsw i64 4, %126
  %128 = add nsw i64 %127, 0
  %129 = getelementptr inbounds double, double* %119, i64 %128
  %130 = load double, double* %129, align 8
  %131 = fmul double %117, %130
  br label %133

132:                                              ; preds = %57
  br label %133

133:                                              ; preds = %132, %112
  %134 = phi double [ %131, %112 ], [ 0.000000e+00, %132 ]
  %135 = fadd double %109, %134
  %136 = load double, double* %14, align 8
  %137 = fcmp ogt double %136, 0.000000e+00
  br i1 %137, label %138, label %158

138:                                              ; preds = %133
  %139 = load double, double* %13, align 8
  %140 = fsub double 1.000000e+00, %139
  %141 = load double, double* %14, align 8
  %142 = fadd double 0.000000e+00, %141
  %143 = fmul double %140, %142
  %144 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %145 = load double*, double** %144, align 8
  %146 = load i64, i64* %11, align 8
  %147 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %148 = load i64, i64* %147, align 8
  %149 = mul nsw i64 %146, %148
  %150 = load i64, i64* %12, align 8
  %151 = add nsw i64 %150, 1
  %152 = add nsw i64 %149, %151
  %153 = mul nsw i64 4, %152
  %154 = add nsw i64 %153, 0
  %155 = getelementptr inbounds double, double* %145, i64 %154
  %156 = load double, double* %155, align 8
  %157 = fmul double %143, %156
  br label %159

158:                                              ; preds = %133
  br label %159

159:                                              ; preds = %158, %138
  %160 = phi double [ %157, %138 ], [ 0.000000e+00, %158 ]
  %161 = fadd double %135, %160
  %162 = load double, double* %13, align 8
  %163 = load double, double* %14, align 8
  %164 = fmul double %162, %163
  %165 = fcmp ogt double %164, 0.000000e+00
  br i1 %165, label %166, label %187

166:                                              ; preds = %159
  %167 = load double, double* %13, align 8
  %168 = fadd double 0.000000e+00, %167
  %169 = load double, double* %14, align 8
  %170 = fadd double 0.000000e+00, %169
  %171 = fmul double %168, %170
  %172 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %173 = load double*, double** %172, align 8
  %174 = load i64, i64* %11, align 8
  %175 = add nsw i64 %174, 1
  %176 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %177 = load i64, i64* %176, align 8
  %178 = mul nsw i64 %175, %177
  %179 = load i64, i64* %12, align 8
  %180 = add nsw i64 %179, 1
  %181 = add nsw i64 %178, %180
  %182 = mul nsw i64 4, %181
  %183 = add nsw i64 %182, 0
  %184 = getelementptr inbounds double, double* %173, i64 %183
  %185 = load double, double* %184, align 8
  %186 = fmul double %171, %185
  br label %188

187:                                              ; preds = %159
  br label %188

188:                                              ; preds = %187, %166
  %189 = phi double [ %186, %166 ], [ 0.000000e+00, %187 ]
  %190 = fadd double %161, %189
  store double %190, double* %15, align 8
  %191 = load double, double* %13, align 8
  %192 = fsub double 1.000000e+00, %191
  %193 = load double, double* %14, align 8
  %194 = fsub double 1.000000e+00, %193
  %195 = fmul double %192, %194
  %196 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %197 = load double*, double** %196, align 8
  %198 = load i64, i64* %11, align 8
  %199 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %200 = load i64, i64* %199, align 8
  %201 = mul nsw i64 %198, %200
  %202 = load i64, i64* %12, align 8
  %203 = add nsw i64 %201, %202
  %204 = mul nsw i64 4, %203
  %205 = add nsw i64 %204, 1
  %206 = getelementptr inbounds double, double* %197, i64 %205
  %207 = load double, double* %206, align 8
  %208 = fmul double %195, %207
  %209 = load double, double* %13, align 8
  %210 = fcmp ogt double %209, 0.000000e+00
  br i1 %210, label %211, label %231

211:                                              ; preds = %188
  %212 = load double, double* %13, align 8
  %213 = fadd double 0.000000e+00, %212
  %214 = load double, double* %14, align 8
  %215 = fsub double 1.000000e+00, %214
  %216 = fmul double %213, %215
  %217 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %218 = load double*, double** %217, align 8
  %219 = load i64, i64* %11, align 8
  %220 = add nsw i64 %219, 1
  %221 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %222 = load i64, i64* %221, align 8
  %223 = mul nsw i64 %220, %222
  %224 = load i64, i64* %12, align 8
  %225 = add nsw i64 %223, %224
  %226 = mul nsw i64 4, %225
  %227 = add nsw i64 %226, 1
  %228 = getelementptr inbounds double, double* %218, i64 %227
  %229 = load double, double* %228, align 8
  %230 = fmul double %216, %229
  br label %232

231:                                              ; preds = %188
  br label %232

232:                                              ; preds = %231, %211
  %233 = phi double [ %230, %211 ], [ 0.000000e+00, %231 ]
  %234 = fadd double %208, %233
  %235 = load double, double* %14, align 8
  %236 = fcmp ogt double %235, 0.000000e+00
  br i1 %236, label %237, label %257

237:                                              ; preds = %232
  %238 = load double, double* %13, align 8
  %239 = fsub double 1.000000e+00, %238
  %240 = load double, double* %14, align 8
  %241 = fadd double 0.000000e+00, %240
  %242 = fmul double %239, %241
  %243 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %244 = load double*, double** %243, align 8
  %245 = load i64, i64* %11, align 8
  %246 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %247 = load i64, i64* %246, align 8
  %248 = mul nsw i64 %245, %247
  %249 = load i64, i64* %12, align 8
  %250 = add nsw i64 %249, 1
  %251 = add nsw i64 %248, %250
  %252 = mul nsw i64 4, %251
  %253 = add nsw i64 %252, 1
  %254 = getelementptr inbounds double, double* %244, i64 %253
  %255 = load double, double* %254, align 8
  %256 = fmul double %242, %255
  br label %258

257:                                              ; preds = %232
  br label %258

258:                                              ; preds = %257, %237
  %259 = phi double [ %256, %237 ], [ 0.000000e+00, %257 ]
  %260 = fadd double %234, %259
  %261 = load double, double* %13, align 8
  %262 = load double, double* %14, align 8
  %263 = fmul double %261, %262
  %264 = fcmp ogt double %263, 0.000000e+00
  br i1 %264, label %265, label %286

265:                                              ; preds = %258
  %266 = load double, double* %13, align 8
  %267 = fadd double 0.000000e+00, %266
  %268 = load double, double* %14, align 8
  %269 = fadd double 0.000000e+00, %268
  %270 = fmul double %267, %269
  %271 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %272 = load double*, double** %271, align 8
  %273 = load i64, i64* %11, align 8
  %274 = add nsw i64 %273, 1
  %275 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %276 = load i64, i64* %275, align 8
  %277 = mul nsw i64 %274, %276
  %278 = load i64, i64* %12, align 8
  %279 = add nsw i64 %278, 1
  %280 = add nsw i64 %277, %279
  %281 = mul nsw i64 4, %280
  %282 = add nsw i64 %281, 1
  %283 = getelementptr inbounds double, double* %272, i64 %282
  %284 = load double, double* %283, align 8
  %285 = fmul double %270, %284
  br label %287

286:                                              ; preds = %258
  br label %287

287:                                              ; preds = %286, %265
  %288 = phi double [ %285, %265 ], [ 0.000000e+00, %286 ]
  %289 = fadd double %260, %288
  store double %289, double* %16, align 8
  %290 = load double, double* %13, align 8
  %291 = fsub double 1.000000e+00, %290
  %292 = load double, double* %14, align 8
  %293 = fsub double 1.000000e+00, %292
  %294 = fmul double %291, %293
  %295 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %296 = load double*, double** %295, align 8
  %297 = load i64, i64* %11, align 8
  %298 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %299 = load i64, i64* %298, align 8
  %300 = mul nsw i64 %297, %299
  %301 = load i64, i64* %12, align 8
  %302 = add nsw i64 %300, %301
  %303 = mul nsw i64 4, %302
  %304 = add nsw i64 %303, 2
  %305 = getelementptr inbounds double, double* %296, i64 %304
  %306 = load double, double* %305, align 8
  %307 = fmul double %294, %306
  %308 = load double, double* %13, align 8
  %309 = fcmp ogt double %308, 0.000000e+00
  br i1 %309, label %310, label %330

310:                                              ; preds = %287
  %311 = load double, double* %13, align 8
  %312 = fadd double 0.000000e+00, %311
  %313 = load double, double* %14, align 8
  %314 = fsub double 1.000000e+00, %313
  %315 = fmul double %312, %314
  %316 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %317 = load double*, double** %316, align 8
  %318 = load i64, i64* %11, align 8
  %319 = add nsw i64 %318, 1
  %320 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %321 = load i64, i64* %320, align 8
  %322 = mul nsw i64 %319, %321
  %323 = load i64, i64* %12, align 8
  %324 = add nsw i64 %322, %323
  %325 = mul nsw i64 4, %324
  %326 = add nsw i64 %325, 2
  %327 = getelementptr inbounds double, double* %317, i64 %326
  %328 = load double, double* %327, align 8
  %329 = fmul double %315, %328
  br label %331

330:                                              ; preds = %287
  br label %331

331:                                              ; preds = %330, %310
  %332 = phi double [ %329, %310 ], [ 0.000000e+00, %330 ]
  %333 = fadd double %307, %332
  %334 = load double, double* %14, align 8
  %335 = fcmp ogt double %334, 0.000000e+00
  br i1 %335, label %336, label %356

336:                                              ; preds = %331
  %337 = load double, double* %13, align 8
  %338 = fsub double 1.000000e+00, %337
  %339 = load double, double* %14, align 8
  %340 = fadd double 0.000000e+00, %339
  %341 = fmul double %338, %340
  %342 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %343 = load double*, double** %342, align 8
  %344 = load i64, i64* %11, align 8
  %345 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %346 = load i64, i64* %345, align 8
  %347 = mul nsw i64 %344, %346
  %348 = load i64, i64* %12, align 8
  %349 = add nsw i64 %348, 1
  %350 = add nsw i64 %347, %349
  %351 = mul nsw i64 4, %350
  %352 = add nsw i64 %351, 2
  %353 = getelementptr inbounds double, double* %343, i64 %352
  %354 = load double, double* %353, align 8
  %355 = fmul double %341, %354
  br label %357

356:                                              ; preds = %331
  br label %357

357:                                              ; preds = %356, %336
  %358 = phi double [ %355, %336 ], [ 0.000000e+00, %356 ]
  %359 = fadd double %333, %358
  %360 = load double, double* %13, align 8
  %361 = load double, double* %14, align 8
  %362 = fmul double %360, %361
  %363 = fcmp ogt double %362, 0.000000e+00
  br i1 %363, label %364, label %385

364:                                              ; preds = %357
  %365 = load double, double* %13, align 8
  %366 = fadd double 0.000000e+00, %365
  %367 = load double, double* %14, align 8
  %368 = fadd double 0.000000e+00, %367
  %369 = fmul double %366, %368
  %370 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %371 = load double*, double** %370, align 8
  %372 = load i64, i64* %11, align 8
  %373 = add nsw i64 %372, 1
  %374 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %375 = load i64, i64* %374, align 8
  %376 = mul nsw i64 %373, %375
  %377 = load i64, i64* %12, align 8
  %378 = add nsw i64 %377, 1
  %379 = add nsw i64 %376, %378
  %380 = mul nsw i64 4, %379
  %381 = add nsw i64 %380, 2
  %382 = getelementptr inbounds double, double* %371, i64 %381
  %383 = load double, double* %382, align 8
  %384 = fmul double %369, %383
  br label %386

385:                                              ; preds = %357
  br label %386

386:                                              ; preds = %385, %364
  %387 = phi double [ %384, %364 ], [ 0.000000e+00, %385 ]
  %388 = fadd double %359, %387
  store double %388, double* %17, align 8
  store double 1.000000e+00, double* %18, align 8
  %389 = load double, double* %15, align 8
  %390 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %391 = load double*, double** %390, align 8
  %392 = load i64, i64* %7, align 8
  %393 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %394 = load i64, i64* %393, align 8
  %395 = mul nsw i64 %392, %394
  %396 = load i64, i64* %8, align 8
  %397 = add nsw i64 %395, %396
  %398 = mul nsw i64 4, %397
  %399 = add nsw i64 %398, 0
  %400 = getelementptr inbounds double, double* %391, i64 %399
  store double %389, double* %400, align 8
  %401 = load double, double* %16, align 8
  %402 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %403 = load double*, double** %402, align 8
  %404 = load i64, i64* %7, align 8
  %405 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %406 = load i64, i64* %405, align 8
  %407 = mul nsw i64 %404, %406
  %408 = load i64, i64* %8, align 8
  %409 = add nsw i64 %407, %408
  %410 = mul nsw i64 4, %409
  %411 = add nsw i64 %410, 1
  %412 = getelementptr inbounds double, double* %403, i64 %411
  store double %401, double* %412, align 8
  %413 = load double, double* %17, align 8
  %414 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %415 = load double*, double** %414, align 8
  %416 = load i64, i64* %7, align 8
  %417 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %418 = load i64, i64* %417, align 8
  %419 = mul nsw i64 %416, %418
  %420 = load i64, i64* %8, align 8
  %421 = add nsw i64 %419, %420
  %422 = mul nsw i64 4, %421
  %423 = add nsw i64 %422, 2
  %424 = getelementptr inbounds double, double* %415, i64 %423
  store double %413, double* %424, align 8
  %425 = load double, double* %18, align 8
  %426 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %427 = load double*, double** %426, align 8
  %428 = load i64, i64* %7, align 8
  %429 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %430 = load i64, i64* %429, align 8
  %431 = mul nsw i64 %428, %430
  %432 = load i64, i64* %8, align 8
  %433 = add nsw i64 %431, %432
  %434 = mul nsw i64 4, %433
  %435 = add nsw i64 %434, 3
  %436 = getelementptr inbounds double, double* %427, i64 %435
  store double %425, double* %436, align 8
  br label %437

437:                                              ; preds = %386
  %438 = load i64, i64* %8, align 8
  %439 = add nsw i64 %438, 1
  store i64 %439, i64* %8, align 8
  br label %52

440:                                              ; preds = %52
  br label %441

441:                                              ; preds = %440
  %442 = load i64, i64* %7, align 8
  %443 = add nsw i64 %442, 1
  store i64 %443, i64* %7, align 8
  br label %46

444:                                              ; preds = %46
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @_resize(%struct.pict* noalias sret %0, %struct.pict* byval(%struct.pict) align 8 %1, i32 %2, i32 %3) #0 {
  %5 = alloca i32, align 4
  %6 = alloca i32, align 4
  store i32 %2, i32* %5, align 4
  store i32 %3, i32* %6, align 4
  %7 = load i32, i32* %5, align 4
  %8 = load i32, i32* %6, align 4
  call void @resize(%struct.pict* sret %0, %struct.pict* byval(%struct.pict) align 8 %1, i32 %7, i32 %8)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @crop(%struct.pict* noalias sret %0, %struct.pict* byval(%struct.pict) align 8 %1, i32 %2, i32 %3, i32 %4, i32 %5) #0 {
  %7 = alloca i32, align 4
  %8 = alloca i32, align 4
  %9 = alloca i32, align 4
  %10 = alloca i32, align 4
  %11 = alloca i64, align 8
  %12 = alloca i64, align 8
  %13 = alloca i64, align 8
  %14 = alloca i64, align 8
  store i32 %2, i32* %7, align 4
  store i32 %3, i32* %8, align 4
  store i32 %4, i32* %9, align 4
  store i32 %5, i32* %10, align 4
  %15 = load i32, i32* %10, align 4
  %16 = load i32, i32* %8, align 4
  %17 = sub nsw i32 %15, %16
  %18 = sext i32 %17 to i64
  %19 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 0
  store i64 %18, i64* %19, align 8
  %20 = load i32, i32* %9, align 4
  %21 = load i32, i32* %7, align 4
  %22 = sub nsw i32 %20, %21
  %23 = sext i32 %22 to i64
  %24 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  store i64 %23, i64* %24, align 8
  %25 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 0
  %26 = load i64, i64* %25, align 8
  %27 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %28 = load i64, i64* %27, align 8
  %29 = mul nsw i64 %26, %28
  %30 = mul nsw i64 %29, 4
  %31 = mul i64 %30, 8
  %32 = call noalias i8* @malloc(i64 %31) #3
  %33 = bitcast i8* %32 to double*
  %34 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  store double* %33, double** %34, align 8
  %35 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %36 = load double*, double** %35, align 8
  %37 = icmp ne double* %36, null
  br i1 %37, label %39, label %38

38:                                               ; preds = %6
  call void @fail_assertion(i8* getelementptr inbounds ([15 x i8], [15 x i8]* @.str.25, i64 0, i64 0))
  br label %39

39:                                               ; preds = %38, %6
  store i64 0, i64* %11, align 8
  br label %40

40:                                               ; preds = %156, %39
  %41 = load i64, i64* %11, align 8
  %42 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %43 = load i64, i64* %42, align 8
  %44 = icmp slt i64 %41, %43
  br i1 %44, label %45, label %159

45:                                               ; preds = %40
  store i64 0, i64* %12, align 8
  br label %46

46:                                               ; preds = %152, %45
  %47 = load i64, i64* %12, align 8
  %48 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 0
  %49 = load i64, i64* %48, align 8
  %50 = icmp slt i64 %47, %49
  br i1 %50, label %51, label %155

51:                                               ; preds = %46
  %52 = load i64, i64* %11, align 8
  %53 = load i32, i32* %8, align 4
  %54 = sext i32 %53 to i64
  %55 = add nsw i64 %52, %54
  store i64 %55, i64* %13, align 8
  %56 = load i64, i64* %12, align 8
  %57 = load i32, i32* %7, align 4
  %58 = sext i32 %57 to i64
  %59 = add nsw i64 %56, %58
  store i64 %59, i64* %14, align 8
  %60 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %61 = load double*, double** %60, align 8
  %62 = load i64, i64* %13, align 8
  %63 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %64 = load i64, i64* %63, align 8
  %65 = mul nsw i64 %62, %64
  %66 = load i64, i64* %14, align 8
  %67 = add nsw i64 %65, %66
  %68 = mul nsw i64 4, %67
  %69 = add nsw i64 %68, 0
  %70 = getelementptr inbounds double, double* %61, i64 %69
  %71 = load double, double* %70, align 8
  %72 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %73 = load double*, double** %72, align 8
  %74 = load i64, i64* %11, align 8
  %75 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %76 = load i64, i64* %75, align 8
  %77 = mul nsw i64 %74, %76
  %78 = load i64, i64* %12, align 8
  %79 = add nsw i64 %77, %78
  %80 = mul nsw i64 4, %79
  %81 = add nsw i64 %80, 0
  %82 = getelementptr inbounds double, double* %73, i64 %81
  store double %71, double* %82, align 8
  %83 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %84 = load double*, double** %83, align 8
  %85 = load i64, i64* %13, align 8
  %86 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %87 = load i64, i64* %86, align 8
  %88 = mul nsw i64 %85, %87
  %89 = load i64, i64* %14, align 8
  %90 = add nsw i64 %88, %89
  %91 = mul nsw i64 4, %90
  %92 = add nsw i64 %91, 1
  %93 = getelementptr inbounds double, double* %84, i64 %92
  %94 = load double, double* %93, align 8
  %95 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %96 = load double*, double** %95, align 8
  %97 = load i64, i64* %11, align 8
  %98 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %99 = load i64, i64* %98, align 8
  %100 = mul nsw i64 %97, %99
  %101 = load i64, i64* %12, align 8
  %102 = add nsw i64 %100, %101
  %103 = mul nsw i64 4, %102
  %104 = add nsw i64 %103, 1
  %105 = getelementptr inbounds double, double* %96, i64 %104
  store double %94, double* %105, align 8
  %106 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %107 = load double*, double** %106, align 8
  %108 = load i64, i64* %13, align 8
  %109 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %110 = load i64, i64* %109, align 8
  %111 = mul nsw i64 %108, %110
  %112 = load i64, i64* %14, align 8
  %113 = add nsw i64 %111, %112
  %114 = mul nsw i64 4, %113
  %115 = add nsw i64 %114, 2
  %116 = getelementptr inbounds double, double* %107, i64 %115
  %117 = load double, double* %116, align 8
  %118 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %119 = load double*, double** %118, align 8
  %120 = load i64, i64* %11, align 8
  %121 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %122 = load i64, i64* %121, align 8
  %123 = mul nsw i64 %120, %122
  %124 = load i64, i64* %12, align 8
  %125 = add nsw i64 %123, %124
  %126 = mul nsw i64 4, %125
  %127 = add nsw i64 %126, 2
  %128 = getelementptr inbounds double, double* %119, i64 %127
  store double %117, double* %128, align 8
  %129 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 2
  %130 = load double*, double** %129, align 8
  %131 = load i64, i64* %13, align 8
  %132 = getelementptr inbounds %struct.pict, %struct.pict* %1, i32 0, i32 1
  %133 = load i64, i64* %132, align 8
  %134 = mul nsw i64 %131, %133
  %135 = load i64, i64* %14, align 8
  %136 = add nsw i64 %134, %135
  %137 = mul nsw i64 4, %136
  %138 = add nsw i64 %137, 3
  %139 = getelementptr inbounds double, double* %130, i64 %138
  %140 = load double, double* %139, align 8
  %141 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %142 = load double*, double** %141, align 8
  %143 = load i64, i64* %11, align 8
  %144 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %145 = load i64, i64* %144, align 8
  %146 = mul nsw i64 %143, %145
  %147 = load i64, i64* %12, align 8
  %148 = add nsw i64 %146, %147
  %149 = mul nsw i64 4, %148
  %150 = add nsw i64 %149, 3
  %151 = getelementptr inbounds double, double* %142, i64 %150
  store double %140, double* %151, align 8
  br label %152

152:                                              ; preds = %51
  %153 = load i64, i64* %12, align 8
  %154 = add nsw i64 %153, 1
  store i64 %154, i64* %12, align 8
  br label %46

155:                                              ; preds = %46
  br label %156

156:                                              ; preds = %155
  %157 = load i64, i64* %11, align 8
  %158 = add nsw i64 %157, 1
  store i64 %158, i64* %11, align 8
  br label %40

159:                                              ; preds = %40
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @_crop(%struct.pict* noalias sret %0, %struct.pict* byval(%struct.pict) align 8 %1, i32 %2, i32 %3, i32 %4, i32 %5) #0 {
  %7 = alloca i32, align 4
  %8 = alloca i32, align 4
  %9 = alloca i32, align 4
  %10 = alloca i32, align 4
  store i32 %2, i32* %7, align 4
  store i32 %3, i32* %8, align 4
  store i32 %4, i32* %9, align 4
  store i32 %5, i32* %10, align 4
  %11 = load i32, i32* %7, align 4
  %12 = load i32, i32* %8, align 4
  %13 = load i32, i32* %9, align 4
  %14 = load i32, i32* %10, align 4
  call void @crop(%struct.pict* sret %0, %struct.pict* byval(%struct.pict) align 8 %1, i32 %11, i32 %12, i32 %13, i32 %14)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @read_image(%struct.pict* noalias sret %0, i8* %1) #0 {
  %3 = alloca i8*, align 8
  store i8* %1, i8** %3, align 8
  %4 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 0
  %5 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %6 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %7 = load i8*, i8** %3, align 8
  call void @_readPNG(i64* %4, i64* %5, double** %6, i8* %7)
  ret void
}

declare dso_local void @_readPNG(i64*, i64*, double**, i8*) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @_read_image(%struct.pict* noalias sret %0, i8* %1) #0 {
  %3 = alloca i8*, align 8
  store i8* %1, i8** %3, align 8
  %4 = load i8*, i8** %3, align 8
  %5 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([19 x i8], [19 x i8]* @.str.28, i64 0, i64 0), i8* %4)
  %6 = load i8*, i8** %3, align 8
  call void @read_image(%struct.pict* sret %0, i8* %6)
  ret void
}

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @write_image(%struct.pict* byval(%struct.pict) align 8 %0, i8* %1) #0 {
  %3 = alloca i8*, align 8
  store i8* %1, i8** %3, align 8
  %4 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 0
  %5 = load i64, i64* %4, align 8
  %6 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 1
  %7 = load i64, i64* %6, align 8
  %8 = getelementptr inbounds %struct.pict, %struct.pict* %0, i32 0, i32 2
  %9 = load double*, double** %8, align 8
  %10 = load i8*, i8** %3, align 8
  call void @_writePNG(i64 %5, i64 %7, double* %9, i8* %10)
  ret void
}

declare dso_local void @_writePNG(i64, i64, double*, i8*) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local void @_write_image(%struct.pict* byval(%struct.pict) align 8 %0, i8* %1) #0 {
  %3 = alloca i8*, align 8
  store i8* %1, i8** %3, align 8
  %4 = load i8*, i8** %3, align 8
  %5 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([18 x i8], [18 x i8]* @.str.29, i64 0, i64 0), i8* %4)
  %6 = load i8*, i8** %3, align 8
  call void @write_image(%struct.pict* byval(%struct.pict) align 8 %0, i8* %6)
  ret void
}

attributes #0 = { noinline nounwind optnone uwtable "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #2 = { noreturn nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #3 = { nounwind }
attributes #4 = { nounwind readnone speculatable willreturn }
attributes #5 = { nounwind readonly "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #6 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #7 = { noreturn nounwind }
attributes #8 = { nounwind readonly }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 10.0.0-4ubuntu1 "}
