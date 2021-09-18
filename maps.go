package opesvengine

import (
    "image"
)

// Return the offset X, Y for a tile given it's relative number
func get_tile_offsets(tile_num int, img_tile_width int) (int, int) {
    var ox, oy int
    ox = tile_num % img_tile_width
    oy = tile_num / img_tile_width
    return ox, oy
}

// Convert a matrix of numbers into a nice map.
// We use this number to find its X and Y on the tilemap file,
// so we get the right subimages. 0 means blank tile
func Draw_map(the_map [][]int, img image.Image,
              frameWidth int, frameHeight int) []Sprite{
    var img_tile_width, ox, oy int
    var tile_map []Sprite
    bounds := img.Bounds()
    img_tile_width = bounds.Max.X / frameWidth
    for j, row := range the_map {
        for i, wall := range row {
            if wall > 0 {
                tile := NewSprite(float64(i * frameWidth),
                                  float64(j * frameHeight),
                                  img, true)
                tile.frame_number = 1
                ox, oy = get_tile_offsets(wall, img_tile_width)
                tile.Frame_offset_x = ox * frameWidth
                tile.Frame_offset_y = oy * frameHeight
                tile_map = append(tile_map, tile)
            }
        }
    }
    return tile_map
}
