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

### Component
1. Drawer: Sprite, Animation, TextBox
2. Size and Position: Drawer's WHXY
3. Cursor: For calculating drawer's relative XY
4. Condition: When drawer to draw 
5. Lifetime: When drawer will not draw 
Put index right after any iterable fields.

### Sprite
1. WH
2. XY
3. Color
4. Tween