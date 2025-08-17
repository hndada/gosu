Decision makings
===
# Confirmed
1. A name should be clear and short
    - Struct name: BacklightComponent -> Backlight
    - field name: cmp -> backlight 
2. Order of setting values: whxy 
3. keysPressed []bool
    - `pressed []bool` seems not a list
4. Prefer position to pos
    - This matches with exported methods (e.g., Position())
5. Avoid wrapping two to three arguments in another struct
    - Resources, Options, and other arguments, explicity.
    - Too dedicated structs harms readability.
6. ChartHeader is passed by value
    - It is about 250 byte, small enough to pass by value 
7. Put `Drawer` suffix when the name conflicts with score-relating terms
    - Combo, Score, Note, Bar, Judgment 

===
# Pending
1. All fields of in-game components are unexported