# Role Playing Game with Computer Opponent using Machine Learning

### The Game

#### Description

This rpg is loosely based off a generic turn based fighting game. The game is
made up of three classes, Knight, Archer, and Wizard. Each class has 6 different
stats: health, stamina (unused), armor, strength, dexterity, and intellect.
Characters die when health drops to or below 0. Armor will prevent character
from taking any damage, but is reduced by the amount of damage done (i.e. a
character with 5 armor that gets attacked for 7 damage will take 0 damage but
will end the turn with 0 armor). Strength, dexterity, and intellect are used for
calculating move success and outcome.

#### Character Generation

When generating a new character, an enemy character is randomly generated and
assigned to that player character. Newly generated characters have randomly
generated attributes based on player class. Stat distribution falls between the
values below:

| Class  | Health | Stamina | Armor | Strength  | Dexterity | Intellect |
|--------|--------|---------|-------|-----------|-----------|-----------|
| Knight | 80-100 | 50-70   | 0-20  | 0.75-1.00 | 0.50-0.75 | 0.25-0.50 |
| Archer | 80-100 | 50-70   | 0-20  | 0.25-0.50 | 0.75-1.00 | 0.50-0.75 |
| Wizard | 80-100 | 50-70   | 0-20  | 0.50-0.75 | 0.25-0.50 | 0.75-1.00 |

#### Game Logic

Each class has the same types of moves. Within the game logic there are three
attacks and three defenses: heavy attack, quick attack, standard attack, block,
parry, and evade. Each attack has a attribute it checks against for success and
a different attribute it checks for damage. This means that a given attack will
have the highest hit rate against enemies whose relevant success attribute is
a low value. Similarly, a given attack will do the most damage when the relevant
damage attribute is high.

| Attack   | Success attribute | Damage attribute |
|----------|-------------------|------------------|
| Heavy    | Intellect         | Strength         |
| Quick    | Strength          | Dexterity        |
| Standard | Dexterity         | Intellect        |

Defensive moves are slightly different. The three moves each rely on a single
attribute to determine success and outcome.

| Defense | Attribute | Success outcome       | Fail outcome       |
|---------|-----------|-----------------------|--------------------|
| Block   | Strength  | Heals character       | Takes enemy damage |
| Parry   | Dexterity | Reflects enemy attack | Takes extra damage |
| Evade   | Intellect | Repairs armor         | Takes enemy damage |

When generating a new character, an enemy character is randomly generated and
assigned to that player character.

### Machine Learning

This portion is all TODO. What follows is dev notes:

#### MinMax Strategy

The AI creates a matrix which represents the minimum and maximum outcomes for
each move. Each cell has an average value associated with it. A positive value
means a net good outcome, and a negative value has a net bad outcome. These
averages are all normalized and the normalized values are used as a weight. A
random number is generated and a move is selected based off of the weighted
values which were generated. This creates a pseudorandom move selection set
where the computer enemy will prefer moves which will have a good chance of
having a positive outcome for itself. Positive outcomes involve damaging the
player, healing its own health, and repairing its armor.

#### Reinforcement Learning

Rewards:
- Small penalty for losing health or armor
- Large penalty for health -> 0
- Small reward for health or armor increased
- Small reward for enemy health or armor decreased
- Large reward for enemy health -> 0

States:
- Health self&&enemy
- Stamina self&&enemy
- Armor self&&enemy
- Strength self&&enemy
- Dexterity self&&enemy
- Intellect self&&enemy

Actions:
- Heavy attack
- Quick attack
- Standard attack
- Block
- Parry
- Dodge

### Awknowledgements

Lili Lommel for character art

Tyler Balson for guidance with AI
