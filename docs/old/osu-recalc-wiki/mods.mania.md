Update mods in osu!mania
===================================================================
Currently the mods HardRock in o!mania just changes parameter of the map, 
which doesnt make a player feel many changes unlike HardRock in other gamemodes. 
I assume HardRock is still 'unrankable' due to that issue. 
This is the suggestion of new HardRock in osu! mania, as well as other mods in mania that can be improved by following suggestion.

## HardRock 
Since OD and SV affects on difficulty, HR now have a chance to be a thing.
The possible features are listed below:
* Multiply amount of deviation of scroll speed change to **1.4x** (which makes SV harsh)
* Make all normal notes to `Short LNs`; this might make a player feel hitting notes like 'staccato'
* Give penalty more on `Stamina` when pausing during the play
* Give penalty on `HitBonusValue` if a player hit wrong keys at empty lane (this is so-called '空POOR'(means 'null MISS') in Beatmania IIDX series)

## Random 
Since now how pattern spread affects the difficulty, the mod `Random`(RD) now changes the SR. 
It is possible that playing RD with formerly generated *seeds* like replay working in 'Lunatic Raves 2'.
But still forcing Random as *unranked* looks decent to me.

## Mirror
Currently the mod `Mirror`(MR) literally flips the beatmap horizontally only.
This makes MR no use in most of 8K maps, since the scratch lane is also flipped to the other lane, 
while it is supposed to be stay in the same lane in BMS.

New MR system will check whether the map is `7+1 keymode` at first, and distinguish a 8K beatmap has *scratch lane* or not.
If the mod system finds the beatmap has scratch lane, it won't flip the lane and make it stayed.  

This goes important in terms of difficulty of a beatmap 
since whether the beatmap has scratch lane or not will change how the difficulty being calculated.