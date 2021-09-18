package opesvengine

import (
    "fmt"
    "bytes"
    "testing"
    "time"
    "image"
    "image/color"
    "image/draw"
    "github.com/hajimehoshi/ebiten/v2" 
    "github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

// draw a simplistic map
var gamemap = [][]int {
            { 0,  0,  0,  0, },
            { 0,  0,  0,  0, },
            { 0,  0,  0,  1, },
            { 1,  1,  1,  1, },
            }

const (
    tileSize = 32 // used to calculate canvas size
    screenWidth  = 4 * tileSize
    screenHeight = 4 * tileSize
)

// Ebiten Game struct
type Game struct {
    World *Universe
    ScreenWidth int
    ScreenHeight int
}

// Ebiten Update loop method
func (g *Game) Update() error {
    GeUpdate(g.World)
    return nil
}

// Ebiten Update draw method
func (g *Game) Draw(screen *ebiten.Image) {
    GeDraw(g.World, screen)
}

// Ebiten Update layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return g.ScreenWidth, g.ScreenHeight
}

// This does the sabe as GeUpdate, but disables keyboard and enemy control
func MockGeUpdate(g *Universe) {
    g.Player.sprite.count++
    for idx := range g.Projectiles {
        g.Projectiles[idx].count++
        // remove projectile when end its max_life
        if g.Projectiles[idx].count > g.Projectiles[idx].max_life {
            g.Projectiles = append(g.Projectiles[:idx], g.Projectiles[idx+1:]...)
            break
        }
    } 
    for idx := range g.Enemies {
        g.Enemies[idx].sprite.count++
    }
    // disable this so the tests control the characters
    // control_enemy(g)
    // grab_controls(g)
    run_physics(g)
}

// Instantiate a game with Map, Player and Enemies
func create_a_game(t *testing.T) *Universe {
    img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
    if err != nil {
        t.Errorf("Error loading image")
    }
    // player character
    player_char := NewCharacter(0, 0, 16, 32, img)
    // draw small square to be used as ammo
    bullet_img := image.NewRGBA(image.Rect(0, 0, 4, 4))
    red := color.RGBA{255, 0, 0, 255}
    draw.Draw(bullet_img, bullet_img.Bounds(), &image.Uniform{red}, image.ZP, draw.Src)
    player_char.Set_weapon(NewWeapon("gun", 5, 39, 10, 200, bullet_img))
    // Enemy
    enemy_char := NewCharacter(0, 0, 16, 32, img)
    
    img, _, err = image.Decode(bytes.NewReader(images.Tiles_png))
    if err != nil {
        t.Errorf("Error loading image")
    }
    g := &Universe{
        Player: player_char,
        Tile_map: Draw_map(gamemap, img, tileSize, tileSize),
        ScreenWidth: screenWidth,
        ScreenHeight: screenHeight,
        MapOffset_x: 0,
        MapOffset_y: 0,
        Enemies: []Character { enemy_char, },
    }
    return g
}

// instantiate a sprite and check if it is ok
func TestSprite(t *testing.T) {
    img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
    if err != nil {
        t.Errorf("Error loading image")
    }
    X := float64(10)
    Y := float64(10)
    sprite := NewSprite(X, Y, img, false)
    
    if sprite.speed_x != 0 {
        t.Errorf("speed different than expected: %f - %d", sprite.speed_x, 0)
    }
    // Check boundaries
    if sprite.rect_top() != Y {
        t.Errorf("rect_top() mismatch: %f - %f", sprite.rect_top(), Y)
    }
    if sprite.rect_bottom() != Y + float64(sprite.Heigth) {
        t.Errorf("rect_top() mismatch: %f - %f", sprite.rect_top(), Y)
    }
    // Bounds
    img_rect := image.Rect(0, 0, 0 + sprite.frame_width, 0 + sprite.frame_heigth)
    sub_img := sprite.eimage.SubImage(img_rect)
    bounds := sub_img.Bounds()
    // Check bounds if they are ok
    if bounds.Max.Y + int(Y) != int(sprite.rect_bottom()) {
        t.Errorf("bounds mismatch: %d - %f", bounds.Max.Y + int(Y), sprite.rect_bottom())
    }
    if bounds.Max.X + int(X) != int(sprite.rect_right()) {
        t.Errorf("bounds mismatch: %d - %f", bounds.Max.X + int(X), sprite.rect_right())
    }
}

// Test gravity and that characters stop at floor
func TestGroundCollision(t *testing.T) {
    var last_speed float64
    g := create_a_game(t)
    // simple test gravity
    MockGeUpdate(g)
    if g.Player.sprite.speed_y <= 0 {
        t.Errorf("gravity is not pulling player down")
    }
    // Test hit the ground
    for !g.Player.sprite.on_the_ground {
        MockGeUpdate(g)
        if g.Player.sprite.speed_y > 0 {
            last_speed = g.Player.sprite.speed_y
        }
    }
    // Check that player reached max speed when falling
    if last_speed < 3.4 {
        t.Errorf("2player is not vertically stoped")
    }
    // And now is resting
    if g.Player.sprite.speed_y != 0 {
        t.Errorf("player is not vertically stoped")
    }
}

// test collision with walls
func TestWallCollision(t *testing.T) {
    var distance int
    g := create_a_game(t)
    // wait for player to hit the ground
    for !g.Player.sprite.on_the_ground {
        MockGeUpdate(g)
    }
    g.Player.sprite.speed_x = 1
    
    distance = 0
    for g.Player.sprite.speed_x > 0 {
        MockGeUpdate(g)
        distance++
    }
    if distance < 2 * tileSize {
        t.Errorf("player didn't walked at all")
    }
    if g.Player.sprite.speed_x != -1 {
        t.Errorf("player didn't bounce at wall")
    }
    g.Player.sprite.speed_x = 0
}

// Test Enemy wall collision
func TestEnemyWallCollision(t *testing.T) {
    var distance int
    g := create_a_game(t)
    // wait for doll to hit the ground
    enemy := g.Enemies[0].sprite
    for !enemy.on_the_ground {
        MockGeUpdate(g)
    }
    // Tell doll to move
    enemy.speed_x = 1
    // She will run untill she hit the wall
    distance = 0
    for enemy.speed_x > 0 {
        MockGeUpdate(g)
        distance++
    }
    if distance < 2 * tileSize {
        t.Errorf("enemy didn't walked at all")
    }
    if enemy.speed_x != -1 {
        t.Errorf("enemy didn't bounced at wall")
    }
    g.Player.sprite.speed_x = 0
}

// test shooting a gun
func TestShootAGun(t *testing.T) {
    var distance int
    g := create_a_game(t)
    // wait for enemy to hit the ground
    enemy := g.Enemies[0].sprite
    for !enemy.on_the_ground {
        MockGeUpdate(g)
    }
    // Move the enemy to a corner
    enemy.speed_x = 1
    // She will run untill she hit the wall
    distance = 0
    for enemy.speed_x > 0 {
        MockGeUpdate(g)
        distance++
    }
    enemy.speed_x = 0
    // Give the player a gun
    img, _, _ := image.Decode(bytes.NewReader(images.Tiles_png))
    g.Player.Set_weapon(NewWeapon("gun", 5, 39, 10, 0, img))
    // Turn player in direction of enemy
    if g.Player.sprite.pos_x < enemy.pos_x {
        g.Player.sprite.direction = RIGHT
    } else {
        g.Player.sprite.direction = LEFT
    }
    // Shoot a gun
    time.Sleep(100 * time.Millisecond)
    g.Player.Shoot(g)
    MockGeUpdate(g)
    if len(g.Projectiles) < 1 {
        t.Errorf("No projectiles found")
    }
    // Wait for the bullet to hit the target
    for g.Projectiles[0].speed_x > 0 {
        MockGeUpdate(g)
    }
    MockGeUpdate(g)
    if g.Enemies[0].health != 100 - g.Player.weapon.damage {
        t.Errorf("The shot didn't hurt %d - %d", g.Enemies[0].health, g.Player.weapon.damage)
    }
}

// Integration test instantiating a game
func TestUniversEngine(t *testing.T) {
    // Player
    img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
    if err != nil {
        t.Errorf("Error loading image")
    }
    player_char := NewCharacter(0, 0, 16, 32, img)
    // Enemy
    img, _, err = image.Decode(bytes.NewReader(images.Runner_png))
    if err != nil {
        t.Errorf("Error loading image")
    }
    enemy_char := NewCharacter(100, 20, 16, 32, img)
    // the Universe object
    img, _, err = image.Decode(bytes.NewReader(images.Tiles_png))
    if err != nil {
        t.Errorf("Error loading image")
    }
    universe := Universe{
        Player: player_char,
        Tile_map: Draw_map(gamemap, img, tileSize, tileSize),
        ScreenWidth: screenWidth,
        ScreenHeight: screenHeight,
        MapOffset_x: 10,
        MapOffset_y: 0,
        Enemies: []Character { enemy_char, },
    }
    // the ebiten Game object
    eg := &Game{
        World: &universe,
        ScreenWidth: screenWidth,
        ScreenHeight: screenHeight,
    }
    fmt.Println("Testing game engine...")
    eg.Update()
    if len(universe.Enemies) < 1 {
        t.Errorf("No enemies found")
    }
}
