Some of the popular properties of `pp maps` is `simple dense chord` and `short LNs`.
On the other hand, there are 'underrated' maps which difficulties are underrated.
'SV' and 'Delay' patterns are example of factors making a map underrated.

# Strain
`Strain` indicates the amount of sheer strain that fingers would feel when playing a beatmap, which star rating is heavily based on.  
To deal with it, here I set the *strain axioms*:
1. Pressing with outer fingers is hard.
2. Pressing with same fingers soon is hard. 
3. Adjacent fingers tend to move along. 
4. Keeping fingers pressing is hard.

Then I assumed the following *strain rules*:
1. Each finger has different base strain value; more outer fingers have higher strain.
2. Pattern which makes a player `jack` (hitting next note at same lane soon) gives more strain value.
3. Pattern which makes a player `trill` (hitting next note at adjacent lanes soon) gives more strain value. 
4. Pattern which makes a player `chord` (hitting next note at adjacent lanes with) gives *less* strain value. 
5. `Hold notes` gives other notes strain bonus.

## Base strain value
The following table infers the current setting of strain of each finger.

idx | finger | base strain | base strain in specific condition
--- | --- | --- | ---
0 | Thumb  | 1.15 | 1.25 (where a keymode forces thumbs not to use a spacebar e.g., 10K)
1 | Index  | 1.0 | -
2 | Middle | 1.1 | -
3 | Ring   | 1.2 | -
4 | Little | 1.3 | 1.15 (at so-called `Scratch-map`)

## Jack
Setting OD as a parameter, the strain multiplier gets higher if the sooner to hit next note in a same lane. 
Min/max value of multiplier is **1** and **2.5** each.

## Trill and Chord
동시에 쳤을 때 miss나는 만큼 차이가 있을때 트릴 보정 최대
동시에 쳤을 때 300이 나기 시작하면, 가까울수록 동시타 보정; strain 패널티
0.75~1.4
2, n개 떨어진 것도 반영해야겠다

기본적으로 OD가 높을수록 시간 여유롭게?
OD가 높을수록 유리할 수도 있을까? miss 판정 범위가 다를테니까
이빨 빠진 동시타로 Test할 필요가 있음. -> 테스트 완료, OD 클수록 miss 판정 범위도 같이 줄어들음.

- 계산에 따라서, 6동시타가 7동시타보다 strain 클 수 있다. 이때, 6키 연타인데 7키 연타처럼 친다면? 
-> 7키 동시타처럼 쳤을 때, 200 이상의 판정이 나온다면 min 값 취하기.
- 홀수키에 대해서, 엄지는 쉬운 손으로 판정; min()
- 트릴 그래프에서 동시타 패널티도 동시에 적용
- 그리고 각 보정은 양측 라인에 중복해서 거는걸로.

## 메모
- (o) 붙어 있는 동시타에 대하여, 6키보다는 7키가 더 strain이 크게 계산되게 할 것. 많이 증가되진 않게. 
- 3계단같이 트릴이 여러 Lane에 걸쳐 있을 경우. (처리 잘 될 것으로 소박히 기대)
- 같은 딜레이여도, 완전 계단은 쉽고 scattered 는 어렵다

## Hold notes (Long notes)
`Hold note`, or `long note`(LN) requires two actions: start to press and pressing it off when it's about to end.

### Strain bonus to normal notes by LNs
For notes while holding LNs gets strain bonus. 
Their strain will be multiplied with multiplier, which starts with 1.0x:
- if outer LN exists: gain **+0.1x** on the multiplier 
- if inner LN exists:
    * if the non-middle inner LN is next to the note : gain **+0.15x** on the multiplier
    * else: gain **+0.1x** on the multiplier
- if other LNs that didn't affect on the multiplier yet exist: gain **+0.03x** on the multiplier per rest LNs 

현재는 코드의 간결성을 위해
remain 에만 base strain 및 impact적용
base는 건들지 않다
모든 롱놋 건드려봤자 상수가 변하지 않는 이상 큰 변화 없을 것 같음

Notes won't get strain bonus from LNs that have just start.
Notes will get partial strain bonus from LNs that have been about to finish or just finished.
Cancel the strain bonus if relating long note is off; regain the strain bonus if the long note is pressed again.

### Strain of LN itself
There are two types of LNs: normal LNs and `short LN`s. 
`Short LN` infers LN which length is short enough that can be considered as a normal note.
Normal LNs would get strain at finish part, while `short LN` wouldn't.
`Short LN` would get a bit more strain at start part, while normal LNs would get base strain value only.

LN types | start part | finish part
--- | --- | ---
normal LNs | `base_strain` | around **0.4** 
short LNs | **1.1**x`base_strain` | **0**


# Stamina
Trace of *Strain*. 
맵 시작 시간부터 맵이 끝날 때까지 반영
stamina=0.9*stamina+0.12(strain at current unit time)
계수 합이 1이 넘는게 point다.
(strain이 지나치게 낮으면 (일시정지 등) 추가 panelty 도입 가능성)
일시정지 패널티 엄격히 하는 것: NF FC를 노모드 FC보다 덜 쳐주는 것과 같은 양상

멈출 경우, 이전 파트의 "stamina" 파트가 시간에 따라 감소
stamina 측정에 그만큼 zero 값이 반영.
이에 Adjusted Score 역시 자연스럽게 감소
그러나, 이를 또 악용하여 나머지 파트의 weight 가 올라가는 일이 없도록
scale 이 역으로 증가함 없이 감소만 하도록 조정.

점수 혹은 pp 감소가 발생.
SR 감소가 있게 할지 안하게 할지는 아직 미정, 그러나 안하게 할듯. 
단, 풀콤 메달 인정 X

stamina, 실시간 측정.

- 일시정지시 모든 손가락의 stamina '밀도 0+alpha' 수준으로 감소 (원하는 만큼 쉴 수 있으니까) 
- 고의 여부는 안따짐, 실수로 멈췄든 일부러 멈췄던 휴식은 취하게 됨 (플레이중 키보드가 고장났다고 해서 보정해주지 않는것과 같은 이치)
- 정 실수로 누르는 경우를 방지하고 싶으면, ESC 세번 이상 눌러야 멈추게 만들면 됨
- 안쳐서 미스면 감소
- 여차하면 일시정지 걸 수 있는 시점에선 '일시정지를 충분히 한다'고 가정 -> stamina zeroize

# Reading
(normal notes는, strain이 reading을 압도한다고 생각하여 고려 안함)
(LN 역시 잠시 고려하였으나, 마찬가지로 strain이 대부분 클 것으로 예상하여 고려 안함)
이 부분에서 다루는 건 SV.

같은 패턴에 대해, SV가 어떻게 주어지냐에 따라 최대 1.08 ~ 1.12 or 1.19(HDFL) 배 부여.

노트 당 노출시간의 비율 * 이전 노트들 노출시간의 비율 경향 (마지막 노트 100개)

리인카네이션같이 저속 시작은 어떡하지
사실 저속 자체가 어려운건데.
"main 속도"를 정하는게 좋을거 같음

윗부분, 중간부분, 끝부분 노출시간 비율 경향 가중치 두기
(내려오면서도 발생하는 변속을 위해)

