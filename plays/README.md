# Beware
* Slice is often tricky
    * Just assigning slice will shallow copy.
    * When ranging over slice, setting value to copied element has no effect.
* All structs and variables in cmd/* package should be unexported.

# Objective
* Manage UI cmponents with each own struct.

# Naming Convention
* Prefix keys- is put when the type is a slice with the length of keyCount.
* Prefix k- for local variables of keys-.
    * Put suffix 'List' when suffix 's' is not available.

* Index variable ends with -i suffix.
* It is fine to put plural suffix (-s) in the middle if putting the suffix at the last is not proper (e.g., adjective)
* Suffix -t for local variable name of -Type

* Avoid using abbreviation in struct name and field name 
unless the name is explicitly supposed to be expressed in abbreviated form.  
    * Exception: img, anim, whxy

* Local variables are encouraged to be written in abbreviated form, which is up to 3 letters. 
    * field name: sprites, anims
    * local name: s, a

* NewXxx returns struct, while LoadXxx doesn't.

# Keyword
ratio: two quantities with the same units
rate: two quantities with the different units

# Code
* When there are local variable and receiver field which points identical value, use local variable: prefer to use shorter one. 

* The convention to organize types and structs in a file by defining the dependencies first and the types that utilize those dependencies later.
    * However, if a small function is used only in a single large function, then put the small at the below of the large.

* If you find any helpful information, commenting is preferred. Grammar check is not necessary at this stage. 

* No XxxArgs. It just makes the code too verbose.
* Introducing interface as a field would make the code too verbose.

* Ratio (speed) first, then time difference.

## Field order
User would look multiple struct's fields in a time. Hence, putting common fields first would be more readable.

### Sprite
1. WH
2. XY
3. Color
4. Tween

### Options
1. The number of drawers (e.g., key count)
2. Arrangement (e.g., order)
3. Sprite options

### Component
1. The number of drawers (e.g., key count)
2. Drawer (e.g., Sprite, Animation)
3. Source of drawer's size and position (e.g., notes)
4. Index of sources: Put index right after any iterable fields
5. Reference point for drawer's position (e.g., cursor)
6. Drawing condition (e.g., keysPressed)
7. Drawer lifetime (e.g., tween, min duration) 

