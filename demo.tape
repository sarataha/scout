Output demo.gif

Set FontSize 20
Set Width 1200
Set Height 600
Set Padding 20
Set WindowBar Colorful
Set Theme "Dracula"

# Start Scout
Type "./scout"
Enter
Sleep 1.5s

# 1. Show Help
Type "?"
Sleep 3s
Type "?"
Sleep 1s

# 2. Cycle Themes (Classic Amber is default)
Type "t"
Sleep 1.5s
Type "t"
Sleep 1.5s
Type "t"
Sleep 1.5s
Type "t"
Sleep 1.5s

# 3. Navigate & Preview
Type "jj"
Sleep 500ms
Type "l"
Enter
Sleep 1s
# Preview cmd/scout/main.go
Type "j"
Sleep 1.5s
Type "h"
Sleep 500ms

# Explore internal/ui
Type "jjj"
Sleep 500ms
Type "l"
Sleep 500ms
Type "jjj"
Sleep 500ms
Type "l"
Sleep 1s
# Preview themes.go
Type "jjjj"
Sleep 2s
Type "h"
Sleep 500ms
Type "h"
Sleep 500ms

# 4. Show Status Line
# Select the binary and use 'o' to open safely (background)
Type "gg"
Type "jjj"
Type "o"
Sleep 2.5s

# Exit
Type "q"
