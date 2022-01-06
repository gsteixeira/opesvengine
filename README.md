OPESV Engine
=============

Orthogonal Perspective Elevation Side View Game Engine.

A simplistic 2d game plataform shooter boiler plate.

OpesV-Engine aims to provide a simplistic platform for 2d games. 

This is **not** production ready and probably never will. This is *Work in Progress*, plus it's done only for enterteinment purposes, so probably it wont be useful to you... yet.

It uses [ebiten](https://ebiten.org/) game library. The images  in the demo game belongs to that project, it is their default image samples.

What can be done with it for now:

    - create a tile based map.
    - control a player that jumps and shoot his gun.
    - control enemies that fiercely tries to kill our hero.
    - collision detection.
    - guns can be changed, and can be parametrized (damage, rate of fire, range, etc)

To be done:

    - looking for a performant way to flip images when player changes direction.
    - allow maps larger than screen.
    - make usage of different tile sizes easier.
    - better AI for npcs.
    - different actions for NPCs (now they just care about one thing: Kill you).
    - a score panel.
    - (SOLVED) still get stuck on the walls sometimes...
    - documentation.

Running the example:
```bash
    git clone https://github.com/gsteixeira/opesvengine
    cd opesvengine/
    # get dependencies
    go mod tidy
    go run sample/demo.go
```

Create a sprite (static object):
```go
    // you can specify an image or use the images described in ebiten
    img := opesvengine.Get_img_from_file("some_image.png")
    // create the sprite
    foo := opesvengine.NewSprite(10, 10, img, false)
```

Create a Character (a live, movable object):
```go
    // specify an image
    img := opesvengine.Get_img_from_file("some_image.png")
    // set initial parameters
    pos_x := 0
    pos_y := 0
    width := 16
    height := 32
    // create the character
    player := opesvengine.NewCharacter(pos_x, pos_y, width, height, img)
```
