package opesvengine

import (
    "math"
    "github.com/hajimehoshi/ebiten/v2"
)

// Grab keyboard controls
func grab_controls (g *Universe) {
    if g.Player.sprite.on_the_ground {
        if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
            g.Player.sprite.speed_x = -g.Player.sprite.grip_x
            g.Player.sprite.direction = LEFT
        } else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
            g.Player.sprite.speed_x = g.Player.sprite.grip_x
            g.Player.sprite.direction = RIGHT
        } else {
            g.Player.sprite.speed_x = 0
        }
        if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
            g.Player.sprite.speed_y = -g.Player.sprite.grip_y
        }
    } else {
        if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
            g.Player.sprite.direction = LEFT
        } else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
            g.Player.sprite.direction = RIGHT
        }
    }
    if ebiten.IsKeyPressed(ebiten.KeySpace) {
        g.Player.Shoot(g)
    }
}

// Control the enemy character
func control_enemy (g *Universe) {
    var distance_x, distance_y float64
    for idx := range g.Enemies {
        // TODO : add other actions like: resting, fighting, patrol.. etc.. 
        // For now, our enemy only cares about one think: to KILL
        // Find out a way to know that we are stuck, if we do, do the oposite thing
        if g.Enemies[idx].sprite.on_the_ground {
            distance_x = math.Abs(
                math.Abs(g.Player.sprite.pos_x) - math.Abs(g.Enemies[idx].sprite.pos_x))
            if distance_x < 50 {
                // dont get too close
                g.Enemies[idx].sprite.speed_x = 0
            } else if g.Player.sprite.pos_x < g.Enemies[idx].sprite.pos_x {
                g.Enemies[idx].sprite.speed_x = LEFT
            } else {
                g.Enemies[idx].sprite.speed_x = RIGHT
            }
            if g.Player.sprite.rect_bottom() < g.Enemies[idx].sprite.pos_y {
                g.Enemies[idx].sprite.speed_y = -g.Player.sprite.grip_y
            }
        }
        // Aim at the enemy
        if g.Player.sprite.pos_x < g.Enemies[idx].sprite.pos_x {
            g.Enemies[idx].sprite.direction = LEFT
        } else {
            g.Enemies[idx].sprite.direction = RIGHT
        }
        // Should shoot or not
        distance_y = math.Abs(
            math.Abs(g.Player.sprite.pos_y) - math.Abs(g.Enemies[idx].sprite.pos_y))
        if distance_y < 20 {
            g.Enemies[idx].Shoot(g)
        }
    }
}
