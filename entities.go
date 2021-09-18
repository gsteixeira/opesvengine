package opesvengine

import (
    "time"
    "image"
)

// The character represents a live entity.
type Character struct {
    sprite *Sprite
    weapon *Weapon
    health int
    score int
}

// Character entity constructor method
func NewCharacter(x, y float64, width, height int, img image.Image) Character {
    sprite := NewSprite(x, y, img, false)
    sprite.Set_dimensions(width, height)
    c := Character{
            sprite: &sprite,
            weapon: nil,
            health: 100,
            score: 0}
    sprite.Parent = &c
    sprite.kind = "character"
    return c
}

// Change character's weapon
func (c *Character) Set_weapon(gun Weapon) {
    c.weapon = &gun
}

// Act of shooting a gun
func (c *Character) Shoot(universe *Universe, ) {
    if c.weapon == nil {
        return
    }
    var fire_pos_x float64
    elapsed := time.Since(c.weapon.last_shot)
    if elapsed.Milliseconds() > c.weapon.fire_rate {
        // Set the fire position at 30% of character height.
        fire_pos_y := c.sprite.pos_y + (float64(c.sprite.Heigth) * 0.30)
        if c.sprite.direction > 0 {
            fire_pos_x = c.sprite.rect_right()
        } else {
            fire_pos_x = c.sprite.pos_x
        }
        bullet := NewSprite(fire_pos_x, fire_pos_y, c.weapon.img, true)
        bullet.kind = "bullet"
        bullet.Parent = c
        bullet.payload = c.weapon.damage
        bullet.gravity = 0 // bullets don't fall down

        bullet.Set_frames(2, 2, 1)
        bullet.Set_dimensions(2, 2)
        bullet.speed_x = c.weapon.max_speed * float64(c.sprite.direction)
        bullet.max_life = int(c.weapon.max_distance)
        universe.Projectiles = append(universe.Projectiles, bullet)
        c.weapon.last_shot = time.Now()
    }
}

// Get character's health - getter
func (c *Character) Get_health() int {
    return c.health
}

// this class represent a weapon and it's attributes
type Weapon struct {
    name string
    max_speed float64
    max_distance float64
    damage int
    // miliseconds between each shot
    fire_rate int64
    last_shot time.Time
    img image.Image
}

// Weapon constructor method
func NewWeapon(name string, max_speed float64,
               max_distance float64, damage int,
               fire_rate int64, img image.Image) Weapon {
    w := Weapon{name, max_speed, max_distance, damage,
                fire_rate, time.Now(), img}
    return w
}
