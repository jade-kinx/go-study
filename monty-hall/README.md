# 몬티홀 문제 (Monty Hall Problem)

[![Monty Hall Problem](http://img.youtube.com/vi/AXB6r-hjsig/0.jpg)](https://www.youtube.com/watch?v=AXB6r-hjsig?t=0s)


몬티 홀(Monty Hall)이라는 캐나다-미국 TV 프로그램 사회자가 진행하던 미국 오락 프로그램 `《Let's Make a Deal》`에서 유래한 `확률` 문제. 아래의 원문은 해당 칼럼에 실린 문제를 그대로 가져온 것이며, 상품의 종류 등의 디테일은 문제에 따라 조금씩 바뀌지만 당연히 수학적인 의미는 동일하다.  

> Suppose you’re on a game show, and you’re given the choice of three doors. Behind one door is a car, behind the others, goats. You pick a door, say #1, and the host, who knows what’s behind the doors, opens another door, say #3, which has a goat. He says to you, "Do you want to pick door #2?" Is it to your advantage to switch your choice of doors?
> 
> 당신이 한 게임 쇼에서 3개의 문 중에 하나를 고를 수 있는 상황이라고 가정하자. 한 문 뒤에는 자동차가, 다른 두 문 뒤에는 염소가 있다. 당신이 1번 문을 고르자, 문 뒤에 무엇이 있는지 아는 사회자는 3번 문을 열어서 염소를 보여줬다. 그리고는 "2번 문으로 바꾸시겠습니까?"라고 물었다. 이 상황에서, 당신의 선택을 바꾸는 게 유리할까?  

세부사항을 쳐내고 핵심적인 규칙만 말하자면 다음과 같다.
* 닫혀 있는 문 3개가 있다.
* 한 문 뒤에는 상품(=자동차)이 있고, 나머지 두 문은 꽝(=염소)이다.
* 참가자는 이 3가지 문 중 하나를 골라야 상품을 얻을 수 있다.
* 참가자가 문 하나를 고르면, 사회자는 남은 2가지 문 중에 하나를 열고 그게 '꽝'이라는 사실을 밝힌다.
* 여기서 참가자에게 다른 문으로 바꿀 수 있는 기회가 주어진다.

원문: [Monty Hall Problem](https://namu.wiki/w/%EB%AA%AC%ED%8B%B0%20%ED%99%80%20%EB%AC%B8%EC%A0%9C)

증명: [let's proove it](./main.go)