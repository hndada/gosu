The issues have happened because the existing algorithm has dealt with the term `Strain` only naively.
The suggesting algorithm will consider `Strain` in more natural way.
The algorithm will think of other new two factors `Stamina` and `Reading`, which help to assess difficulty more accurate.
The algorithm will solve SR under/overrating problems mentioned above.

# Calculate total difficulty  
Calculate `Reading * Strain + Stamina` for every unit.
Sort in descending order, then sum in power series with decay factor; this is same with existing algorithm.
