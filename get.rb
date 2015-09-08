f = ""
(35..36).each do |w|
  f += " hrksnc_#{w}.png"
  files = ""
  (29..30).each do |h|
    files += " hrksnc_#{w}_#{h}.png"
    next if File.exist?("hrksnc_#{w}_#{h}.png")
    `curl -o hrksnc_#{w}_#{h}.png http://www.jma.go.jp/jp/highresorad/highresorad_tile/HRKSNC/201507091240/201507091240/zoom6/#{w}_#{h}.png`
  end
  `convert -append #{files} hrksnc_#{w}.png` unless File.exist?("hrksnc_#{w}.png")
  `rm -f #{files}`
end
`convert +append #{f} hrksnc.png`
`convert -resize 800% hrksnc.png hrksnc_large.png`
`rm -f #{f}`

f = ""
(140..147).each do |w|
  f += " map_#{w}.png"
  files = ""
  (116..123).each do |h|
    files += " map_#{w}_#{h}.png"
    next if File.exist?("map_#{w}_#{h}.png")
    `curl -o map_#{w}_#{h}.png http://www.jma.go.jp/jp/highresorad/map_tile/MAP_COLOR/none/anal/zoom8/#{w}_#{h}.png`
  end
  `convert -append #{files} map_#{w}.png` unless File.exist?("map_#{w}.png")
  `rm -f #{files}`
end
`convert +append #{f} map.png`
`convert -resize 200% map.png map_large.png`
`rm -f #{f}`

f = ""
(280..295).each do |w|
  f += " manucipality_#{w}.png"
  files = ""
  (232..247).each do |h|
    files += " manucipality_#{w}_#{h}.png"
    next if File.exist?("manucipality_#{w}_#{h}.png")
    `curl -o manucipality_#{w}_#{h}.png http://www.jma.go.jp/jp/highresorad/map_tile/MUNICIPALITY/none/none/zoom9/#{w}_#{h}.png`
  end
  `convert -append #{files} manucipality_#{w}.png` unless File.exist?("manucipality_#{w}.png")
  `rm -f #{files}`
end
`convert +append #{f} manucipality.png`
`rm -f #{f}`

f = ""
(280..295).each do |w|
  f += " map_mask_#{w}.png"
  files = ""
  (232..247).each do |h|
    files += " map_mask_#{w}_#{h}.png"
    next if File.exist?("map_mask_#{w}_#{h}.png")
    `curl -o map_mask_#{w}_#{h}.png http://www.jma.go.jp/jp/highresorad/map_tile/MAP_MASK/none/none/zoom9/#{w}_#{h}.png`
  end
  `convert -append #{files} map_mask_#{w}.png` unless File.exist?("map_mask_#{w}.png")
  `rm -f #{files}`
end
`convert +append #{f} map_mask.png`
`rm -f #{f}`

`convert map_large.png hrksnc_large.png -composite tmp1.png`
`convert tmp1.png map_mask.png -composite tmp2.png`
`convert tmp2.png manucipality.png -composite jma.png`
`rm -f tmp1.png tmp2.png`
