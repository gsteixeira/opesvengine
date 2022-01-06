// OPESV Game Engine. 
// Orthographic Projection Elevation Side View Game Engine example.

package main

import (
    "fmt"
    "log"
    "time"
    "bytes"
    "image"
    "image/color"
    "image/draw"
    "math/rand"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/examples/resources/images"
    "github.com/gsteixeira/opesvengine"
)
// map json - tiles matrix, file name, name, 
var gamemap = [][]int {
            { 0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 0,  0,  0,  0,  0},
            { 0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 57, 13, 13, 13, 13},
            { 0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0},
            { 0,  0,  0,  0,  0, 50, 50, 50,  0,  0,  0,  0,  0,  0,  0},
            { 0,  0,  0, 50, 50,  0,  0,  0,  0, 57,  0,  0,  0,  0,  0},
            {57, 57,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0},
            { 0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0, 13, 13, 13},
            { 0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0},
            { 0, 57,  0,  0,  0,  0,  0,  0,  0,  0, 57,  0,  0,  0,  0},
            {50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50,  0,  0,  0},
        }
    
const (
    tileSize     = 32 
    screenWidth  = 15 * tileSize
    screenHeight = 10 * tileSize
)

type Game struct {
    World opesvengine.Universe
    ScreenWidth int
    ScreenHeight int
}

// ebiten Update game loop
func (g *Game) Update() error {
    // remove dead players
    for i := range g.World.Enemies {
        if g.World.Enemies[i].Get_health() <= 0 {
            g.World.Enemies = append(g.World.Enemies[:i], g.World.Enemies[i+1:]...)
        }
    }
    // add a new enemy if the one we have has died
    if len(g.World.Enemies) < 1 {
        g.World.Enemies = append(g.World.Enemies, make_enemy())
    }
    // ebiten update
    opesvengine.GeUpdate(&g.World)
    return nil
}

// the draw function
func (g *Game) Draw(screen *ebiten.Image) {
    opesvengine.GeDraw(&g.World, screen)
}

// the layout function
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return g.ScreenWidth, g.ScreenHeight
}

// Creates a new enemy instance at random location
func make_enemy() opesvengine.Character {
    // draw small square to be used as ammo
    bullet_img := image.NewRGBA(image.Rect(0, 0, 4, 4))
    col := color.RGBA{255, 255, 0, 255}
    draw.Draw(bullet_img, bullet_img.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)
    img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
    if err != nil { panic(err) }
    pos_x := float64(rand.Intn(screenWidth))
    enemy_char := opesvengine.NewCharacter(pos_x, 1, 16, 32, img)
    enemy_char.Set_weapon(opesvengine.NewWeapon("gun", 3, 40, 10, 500, bullet_img))
    return enemy_char
}

func main() {
    // seed random generator
    rand.Seed(time.Now().UnixNano())
    // Create the player
    img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
    player_char := opesvengine.NewCharacter(0, 0, 16, 32, img)
    // draw small square to be used as ammo
    bullet_img := image.NewRGBA(image.Rect(0, 0, 4, 4))
    col := color.RGBA{255, 0, 0, 255}
    draw.Draw(bullet_img, bullet_img.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)
    player_char.Set_weapon(opesvengine.NewWeapon("gun", 3, 40, 10, 200, bullet_img))
    // enemy
    enemies := []opesvengine.Character {make_enemy(), }
    // Map
    img, _, err = image.Decode(bytes.NewReader(images.Tiles_png))
    the_map := opesvengine.Draw_map(gamemap, img, 32, 32)
    if err != nil { panic(err) }
    // the Universe object
    universe := opesvengine.Universe{
        Player: player_char,
        Tile_map: the_map,
        ScreenWidth: screenWidth,
        ScreenHeight: screenHeight,
        MapOffset_x: 10,
        MapOffset_y: 0,
        Enemies: enemies,
    }
    // the ebiten game object
    eg := &Game{
        World: universe,
        ScreenWidth: screenWidth,
        ScreenHeight: screenHeight,
    }
    ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
    ebiten.SetWindowTitle("OPESV GameEngine Demo")
    // Start
    fmt.Println("Starting game...")
    if err := ebiten.RunGame(eg); err != nil {
        log.Fatal(err)
    }
}

