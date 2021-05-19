## SSTable
sparseIndex를 메모리에 유지한다.\
Get 요청이 들어오면 해당 인덱스만큼의 사이즈만 파일에서 불러와 메모리에 탑재한 후, binarySearch로 검색한다. (fseek)
