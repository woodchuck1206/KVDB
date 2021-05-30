# To do
각 Level을 Tier로 나누는 Tiering 방식을 implement하려고 했지만, 비어있는 SSTable level을 어떻게 tier별로 나눌것인가라는 문제가 남아있다.
Level이 threshold에 도달했을때, 이 level은 다음 level의 각 tier에 merge되는 것이 tiering 방식인데, 그렇다면 다음 단계를 키의 range로 미리 tier를 나눠놓아야하는가, 아니면 동적으로 tier의 크기를 조절하는 기능을 넣어야하는가?

## SSTable
sparseIndex를 메모리에 유지한다.\
Get 요청이 들어오면 해당 인덱스만큼의 사이즈만 파일에서 불러와 메모리에 탑재한 후, binarySearch로 검색한다. (fseek)
