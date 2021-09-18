package opesvengine

import (
    "image"
    _ "image/png"
    "log"
    "os"
    "github.com/hajimehoshi/ebiten/v2"
)

// Constants for physics
const (
    GRAVITY = 0.1
    JUMP_POWER = 4
    FWD_SPEED = 1
    LEFT = -1
    RIGHT = 1
    TILESIZE = 32
    FRAME_WIDTH  = TILESIZE
    FRAME_HEIGTH = TILESIZE
)

// Universe holds everything related to the game
type Universe struct {
    Player Character
    Enemies []Character
    Projectiles []Sprite
    Tile_map []Sprite
    Keys []ebiten.Key
    // Screen size
    ScreenWidth int
    ScreenHeight int
    // Tile size in pixels
    TileSize_X int
    TileSize_Y int
    // To be used by camera
    MapOffset_x int
    MapOffset_y int
}

// Ebiten Update method of GamEngine
func GeUpdate(g *Universe){
    g.Player.sprite.count++
    for idx := range g.Projectiles {
        g.Projectiles[idx].count++
        if g.Projectiles[idx].count > g.Projectiles[idx].max_life {
            g.Projectiles = append(g.Projectiles[:idx], g.Projectiles[idx+1:]...)
            break
        }
    }
    for idx := range g.Enemies {
        g.Enemies[idx].sprite.count++
    }
    control_enemy(g)
    grab_controls(g)
    run_physics(g)
}

// Ebiten Draw method of GamEngine Draw map.
func GeDraw(g *Universe, screen *ebiten.Image){
    for _, tile := range g.Tile_map {
            tile.draw(screen)
        }
    // Draw player
    g.Player.sprite.draw(screen)    
    // shoots
    for idx := range g.Projectiles {
        g.Projectiles[idx].draw(screen)
    }
    // Draw enemies
    for idx := range g.Enemies {
        g.Enemies[idx].sprite.draw(screen)
    }
}

// return an image from a local file
func Get_img_from_file(filename string) image.Image {
    f, err := os.Open(filename)
    if err != nil  { panic(err) }
    defer f.Close()
    img, _, err := image.Decode(f)
    if err != nil {
        log.Fatal(err)
    }
    return img
}
