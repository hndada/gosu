# Special Tags 
Special Tags are tags that affect to in-game.
Multiple keys which are semantically same may do the same thing.
Valid value types over a type goes generous: 
`true, 1, on, ...`, or even just key only are all valid values of boolean.
(Its value will be linted in golang standard though: `true`)

# Boolean type tags
Boolean type tags works as `true` when no values are provided.
Prefix `No` or `no` would work as `false` 

## Example
Pitch: `Pitch` value goes true
NoPitch: `Pitch` value goes false

# Semantically same name
NoPitch: PascalCase. standard
noPitch: camelCase
no_pitch: snake_case
nopitch: lowercase

Pitch (nightcore, nc)[`boolean`]: Pitch applied when time rate goes up (highest priority)  

Vocal[`boolean`]: No pitch applied when time rate goes up

Level (lv, diff, difficulty)[`integer`]: Apply custom level value

LevelO2jam (lvo2)[`interger`]

LevelOsu (lvosu, stars, sr)[`float64`]
