넣고 싶은 기능이 많아지니 
'이정도면 내가 그냥 framework 로 하나 짜는게 나을지도 모름' 하면서 개발

Origin
- golang brought me here 
- do*g*
- si kucin*g*/anjin*g*
- 고수 itself

의경가면 osu!대신 gosu! 할거임 

(타 리듬게임에서 쓸모있어 보이는 기능 import)

라이센스 확인
osu!framework - 저작자만 밝히면 상업적으로도 사용 가능
resource - 비상업적 용도로 사용 가능

no converted

4k~10k only

ojm, ojn/bms/sim 파일 불러오기 (low priority so far)

새 sr
새 점수 시스템
- weighted score
- no extra combo tick to ln, only at start and end, 2 in total

mp3 바뀌는거 감지

MariaDB + DBeaver(GUI)

sr (score) / pp 버전 따로 관리


OD/HP 방식 통일
(o2jam, BMS, osu! OD9HP8.5)
miss range 는 od에 영향받지 않게 
(-> 사실 판정 범위 하나로 쓰는 이상 자동이긴 함)


---
곡 선택 시 beatmap 최신버전인지 확인
스펙시 내 pb 보여주기 (그냥 내 기록들을 넘겨주기)

롱놋 일부러 안 누르기 꼼수 해결
기존: 없을 때 기준의 strain으로 계산, 점수 반영 
sr까지 재측정할 수 있겠지만 지속 대미지로 대충 해결 가능일듯
- 안누르고 있는 한 보너스 점수 풀로 회복이 안되게? (안누르는 롱놋 n개당 최대 회복 가능 수치 -10)

---
버전 관리
알파/베타 버전 v0.X.X (diff/pp 안정화 이전 단계)
