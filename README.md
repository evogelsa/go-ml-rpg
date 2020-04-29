# Machine Learning RPG AI -- sp20-id-0007

| Ethan Vogelsang
| evogelsa@iu.edu
| Indiana University
| hid: sp20-id-0007
| lab github: [:cloud:](https://github.iu.edu/ise-engr-e222/sp20-id-0007)
| project repo: [:cloud:](https://github.iu.edu/evogelsa/go-ml-rpg)

---

Keywords: Game, AI, Reinforcement Learning, QTable, Golang 

---

## Abstract

Game designers have many options when it comes to developing computer players
for their games; however, machine learning is rarely considered part of these
options. This may be due to the need for large amounts of data in order to
train a model, or potentially it is due to needing the game to be close to a
deployable state before development can even begin on the AI. Regardless of the
motivations of professional game developers, creating a computer opponent with a
machine learning algorithm is an interesting experiment to explore the limits of
the technology in this field. By using a continuously training algorithm such as
reinforcement learning, the computer opponent will always be learning and
adapting to the strategies players use. Despite some limitations in the design,
the computer opponent proves to be difficult match for most beginner players
and still puts up a fight against those with more experience.

## Introduction

What does machine learning bring to the table that common computer opponents do
not have? Realistically the answer may be that it brings very little, but
in this application it provides a lot of flexibility for a very specific use 
case. This RPG was intentionally designed to work with a computer opponent that
uses reinforcement learning to its advantage. The computer is always learning 
from the games it plays. This allows the computer to continuously be adapting
and adjusting its play to changes in player strategy or game meta. In the ideal
scenario players play against a computer which trains only to their play style,
and thus it is always providing a challenge to the player until they find a new
strategy to adapt.

## Game Design 

#### Description

This rpg is loosely based off a generic turn based fighting game. The game is
made up of three classes, Knight, Archer, and Wizard. Each class has 5 different
stats: health, armor, strength, dexterity, and intellect.  Characters die when
health drops to or below 0. Armor will prevent character from taking any damage,
but is reduced by the amount of damage done (i.e. a character with 5 armor that
gets attacked for 7 damage will take 0 damage but will end the turn with 0
armor). Strength, dexterity, and intellect are used for calculating move
success and outcome.

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
attribute to determine success and outcome. Defensive moves can also only be
successful if the opponent uses an offensive move.

| Defense | Attribute | Success outcome       | Fail outcome       |
|---------|-----------|-----------------------|--------------------|
| Block   | Strength  | Heals character       | Takes enemy damage |
| Parry   | Dexterity | Reflects enemy attack | Takes extra damage |
| Evade   | Intellect | Repairs armor         | Takes enemy damage |

## AI Implementation

The main focus of this experiment is of course to study the viability of machine
learning as a game AI, and in order to do so, a few other AI methods were used
as a measure of comparison. Each method is detailed below.

When running the server there are a number of command line flags.  Running with
`-h` or `--help` will provide a useful text containing the flags available.

#### Random Strategy

The AI picks a random move to use.

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

The reinforcement learning strategy consists of using a QTable to determine
the best option. By default the learning rate is set to .05, discount factor
set to .3, and explore rate set to .05. These can be changed through the run
flags.

The QTable is a table of values which enumerate each possible state and the
potential rewards associated with each action that could be taken from that
state. An example QTable may look similar to the following:

|             | **Attack** | **Defend** |
| **State 1** | 0.52       | 1.50       |
| **State 2** | -5         | 294        |
| **State 3** | 1          | -0.25      |

This QTable is much simpler than the one used in the game, but the concept is
identical. If the computer analyzes the current game and decides the state is
state 1, then it will most likely select the defensive action as that provides
the highest reward. Depending on the exploration rate, there is a chance that
the action is randomly selected rather than selected based off reward. 

The size of the QTable grows exponentially as more states and more actions are
added to the game, which limits the amount of complexity greatly. More states
and more actions means larger file sizes which results in performance loss when
saving and reading data.

Within the game states are determined using health, armor, and class where each
stat is turned into a discreet value between 0 and 2 (inclusive). This gives 729
total states. These descriptors were selected in order to give a general idea of
what the game looks like without increasing complexity by analyzing every
available stat. Possible improvements to the AI difficulty could be found by
increasing the number of states.

Assuming that learning is enabled when the server is run, each time the player
makes an action the game updates its state. The computer receives a reward if
the state changes to the benefit of the AI, and penalties are given if the state
changes to the benefit of the player.

Rewards:
- Player loses health
- Player loses armor
- AI gains health
- AI gains armor

Penalties:
- AI loses health
- AI loses armor

This method of machine learning provides a setup which allows for the computer
to not necessarily need a fully trained model in order to be effective. In fact
during the first round of testing the computer was initially trained against
another computer using randomly generated characters and making decisions with 
the random strategy. This allowed for a relatively efficient way the for 
computer to get experience playing against all possible scenarios. Later into
development human players were recruited to test the viability of the AI.

## Conclusion

The results of player testing showed that the success of the AI was largely
dependent on player experience. Newer players often struggled against the
computer and were unable to defeat it. However, once players were able to
discover a strategy, such as spamming parry as an archer, the AI became much
easier to defeat. From a game design perspective this may actually be a good
thing since the player should be able to feel accomplishment through progress.

During play testing the AI was set to have pretty gentle learning rates, and it
never trained on only one player at a time. This limited the ability for the
algorithm to adapt to the strategies of just one person at a time, and instead
the computer ended up with an average of all players. For the small number of
people playing at once this still turned out to be decently effective, but with
larger player bases it is likely that this implementation would quickly become
ineffective. This could be circumvented by adding an authentication system
allowing players to each have their own AI which has only trained against them.

## Acknowledgements

Special thanks to Lili Lommel for working on the character art which adds a much
needed taste of flavor to the UX.

Thank you to all the play testers who took some time to play the game and help
train the AI and give their feedback.

Please see [report.bib](/blob/master/report.bib) for citations used.
