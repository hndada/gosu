실제 플레이 점수 새 계산 체계 (롱놋 놓을때, 일시정지할때.)

# Effective Score
This aims to solve abusing existing score/pp calculation algorithm.
Certain amount of hard patterns of a beatmap affects most of its star rating, while the distribution of `scores` that each note gets are all same.
So its possible that 'game the algorithm': putting minimum amount of hard patterns on a beatmap and put easy patterns at rest.
This is prone to long 7k maps e.g., [Imperishable Night 2006 [7k Lunatic]](https://osu.ppy.sh/beatmapsets/92190#mania/249346), 
[starlights feat. TEA [Celestial Radiance // 7K]](https://osu.ppy.sh/beatmapsets/831653#mania/1742328).

In suggesting score system, each note gives different score depending on its total strain.
The harder notes will give more score than other easy notes when they are hit.
The max value of score which each note gives is proportional to (strain of the note) / (total strain).
The total score would be **1,000,000** like before.

With new score system, players tend to get lower score on the same beatmap with same skill than before,
since the score is weighted more on harder part, which is hard to get full scores.
To diminish score decrease, the new system will change the ratio on `HitValue` and `HitBonusValue` from **5:5** to **7:3**.

## pp calculation with new score system
Structure of the formula for pp calculation will be same as existing one: 
using a function with taking parameter as `OD` and `SR`.
The parameter `HP` had not effected on pp calculation so far, but new one will try to consider it for natural pp calculation.
Lower score at a beatmap with lower `HP`, higher `SR` will give less pp. (cf. [Doppelganger [Alter Ego]](https://osu.ppy.sh/beatmapsets/407153#mania/884617))

### Dealing with existing score records
There are lots of scores already, and it is essential to calculate *estimated pp* for all of them.
We might be able to design function of score-new pp mapping for each beatmap with using its main difficulty factors. 

* `Stamina` goes **0** whenever total `Strain` of unit pattern belows critical value. (Since it has some chance to abuse pause) 
* Mediocre score at unbalanced map in difficulty gets penalty on pp calculation (cf. [Imperishable Night 2006 [7k Lunatic]](https://osu.ppy.sh/beatmapsets/92190#mania/249346))
    - Might use the ratio on total strain of each tenth fractile at strain order.

## Issue on applying Effective Score system 
The possible problem of Effective Score is that the total score for same play shall be changed 
if the algorithm is updated; strain of each note is newly calculated.
There are some possible solutions for this: 
* A. Stacks minimum requirement trace from all replay.
Strength: One can save actual best performancing play regardless of how algorithms work.  
Weakness: It requires much disk space and intense time and power consumption if algorithm goes updated. 
If the suggesting algorithm is stable and fit well enough, maybe not much change would happen on the scores: less effective comparing to its cost.

* B. No stacks, just let scores based on old algorithm be itself.
Strength: Few burden like as before.
Weakness: Each version of algorithm in scores in the leaderboard goes highly different, which looks confusing to most of people. 
Since the weakness of plan B is fatal, I might consider the A one.
