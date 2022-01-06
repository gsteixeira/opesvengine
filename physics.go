package opesvengine

import (
    "fmt"
    "math"
    "reflect"
)

// Get the collision vector between two sprites
func get_collision_vector(source *Sprite, target Sprite) (float64, float64) {
    source_x, source_y := source.pos_center()
    target_x, target_y := target.pos_center()
    vect_x := target_x - source_x
    vect_y := target_y - source_y
    return vect_x, vect_y
}

// Get the angle of a given vector
func get_vector_angle(vect_x float64, vect_y float64) float64 {
    radians := math.Atan2(vect_y, vect_x)
    degrees := 180 * radians / math.Pi
    return degrees
}

// Tell which side of the object the collision happened
func get_collision_side(source *Sprite, target Sprite) string {
    dist_x, dist_y := get_collision_vector(source, target)
    if math.Abs(dist_x) > math.Abs(dist_y) {
        return "side"
    } else if dist_y > 0 { // moving down
        return "ground"
    } else {
        return "head"
    }
}

// Check collision. Receives a list of targets to check against
func (s *Sprite) check_collision(target_list []Sprite, game *Universe) {
    var horizontal_match bool
    var vertical_match bool
    var side string
    var tile_rect SpriteRect
    var s_rect SpriteRect = s.get_rect(s.speed_x, s.speed_y)
    s.on_the_ground = false
    for idx, tile := range target_list {
        tile_rect = tile.get_rect(0, 0)
        horizontal_match = (s_rect.left < tile_rect.right && s_rect.right > tile_rect.left)
        vertical_match = (s_rect.bottom > tile_rect.top && s_rect.top < tile_rect.bottom)
        if (horizontal_match && vertical_match) {
            side = get_collision_side(s, tile)
            switch side {
                case "side":
                    if !s.on_the_ground {
                        s.speed_x = -s.speed_x
                    } else {
                        s.speed_x = 0
                    }
                case "ground":
                    s.speed_y = 0
                    s.on_the_ground = true
                case "head":
                    if s.speed_y < 0 {
                        s.speed_y = s.gravity * 2
                    }
            }
            // if it's a bullet, inflict damage
            if s.kind == "bullet" && tile.kind == "character" {
                target_list[idx].Parent.health -= s.payload
                s.max_life = 0
                // Look for real reference of the victim to update health
                // There got to be a better way to do this..
                for j := range game.Enemies {
                    if reflect.DeepEqual(game.Enemies[j].sprite, &target_list[idx]) {
                        game.Enemies[j].health = target_list[idx].Parent.health
                        break
                    }
                }
            }
        }
    }
    // Out of screen
    if s.pos_x < 0 {
        s.pos_x = 0
    } else if s.pos_x > float64(game.ScreenWidth) {
        // shoud move to other map, or maybe, you know.. die.
        s.pos_x = float64(game.ScreenWidth)
    }
    if s.pos_y < 0 {
        s.pos_y = 0
        s.speed_y = 0
    } else if s.pos_y > float64(game.ScreenHeight) {
        s.pos_y = 1
        s.pos_x = 1
        fmt.Println("You died!")
    }
}

// Run physics like gravity, movement based on speed, collisions, etc.
func (s *Sprite) run_physics(game *Universe){
    // gravity
    s.speed_y += s.gravity
    // Collisions
    s.check_collision(game.Tile_map, game)
    // Set player's next position
    s.pos_x += float64(s.speed_x)
    s.pos_y += float64(s.speed_y)
}

// Run physics for each moving sprite
func run_physics(u *Universe){
    u.Player.sprite.run_physics(u)
    var enemy_sprites []Sprite
    for idx := range u.Enemies {
        u.Enemies[idx].sprite.run_physics(u)
        enemy_sprites = append(enemy_sprites, *u.Enemies[idx].sprite)
    }
    live_beings := append(enemy_sprites, *u.Player.sprite)
    for idx := range u.Projectiles {
        u.Projectiles[idx].run_physics(u)
        u.Projectiles[idx].check_collision(live_beings, u)
    }
}
