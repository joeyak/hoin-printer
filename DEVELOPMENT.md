General Commands:

- [x] write raw bytes
  - Write()
- [x] read raw bytes
  - Read()
- [x] Print
- [x] Printf
- [x] Println

Programmer Manual Commands:

- [x] HT ~ Horizontal Tab
  - HT()
- [x] LF ~ Print and line feed
  - LF()
- [x] CR ~ Print and carriage return
  - CR()
- [ ] DLE EOT n ~ Real-time status transmission
- [ ] DLE DC4 n m t ~ Generate pulse at real-time
- [ ] ESC SP n ~ Set right-side character spacing
- [ ] ESC ! n ~ Select print mode(s)
- [ ] ESC $ nL nH ~ Set absolute print position
- [ ] ESC % n ~ Select/cancel user-defined character set
- [ ] ESC & y c1 c2 [x1 d1...d(x×x1)]...[xk d1...d(y×xK)] ~ Define user defined characters
- [ ] ESC \* m nL nH d1... dk ~ Select bit-image mode
- [ ] ESC - n ~ Turn underline mode on/off
- [x] ESC 2 ~ Select default line spacing
  - ResetLineSpacing()
- [x] ESC 3 n ~ Set line spacing
  - SetLineSpacing()
  - [ ] Standard Mode
  - [ ] Page mode
- [ ] ESC = n ~ Set peripheral device
- [ ] ESC ? n ~ Cancel user-defined characters
- [X] ESC @ ~ Initialize printer
  - Initialize()
- [x] ESC D n1...nk NUL ~ Set horizontal tab positions
- [ ] ESC E n ~ Turn emphasized mode on/off
- [ ] ESC G n ~ Turn on/off double-strike mode
- [x] ESC J n ~ Print and feed paper
  - Feed()
- [ ] ESC M n ~ Select character font
- [ ] ESC V n ~ Turn 90 degress clockwise rotation mode on/off
- [ ] ESC Z m n k dL dH d1...dn ~ print qr.code
- [ ] ESC \\ nL nH ~ Set relative print position
- [ ] ESC a n ~ Select justification
- [ ] ESC c 3 n (\*) ~ Select paper sensor(s) to output paper end signals
- [ ] ESC c 4 n (\*) ~ Select paper sensor(s) to stop printing
- [ ] ESC C 5 n ~ Enable/disable panel buttons
- [x] ESC d n ~ Print and feed n lines
  - FeedLines()
- [ ] ESC p m t1 t2 ~ Generate pulse
- [ ] ESC t n ~ Select character code table
- [ ] ESC { n ~ Turns on/off upside-down printing mode
- [ ] FS p n m ~ Print NV bit image
- [ ] FS q n [xL xH yL yH d1...dk]<sub>1</sub>...[xL xH yL yH d1...dk]<sub>n</sub> ~ Define NV bit image
- [ ] GS ! n ~ Select character size
- [ ] GS $ nL nH ~ Set absolute vertical print position in page mode
- [ ] GS \* x y d1...d(xxyx8) ~ Define downloaded bit image
- [ ] GS / m ~ Print downloaded bit image
- [ ] GS B n ~ Turn white/black reverse printing mode
- [ ] GS H n ~ Select printing position for HRI characters
- [ ] GS L nL nH ~ Set left margin
- [x] GS V m ~ Select cut mode and cut paper
  - Cut()
- [ ] GS W nL nH ~ Set printing area width
- [ ] GS f n ~ Select font for Human Readable Interpretation (HRI) characters
- [ ] GS h n ~ Select bar code height
- [ ] GS k m d1...dk NUL ~ Print bar code
- [ ] GS k m n d1...dn ~ Print var code
- [ ] GS v 0 m xL xH yL yH d1...dk ~ Print raster bit image
- [ ] GS w n ~ Set bar code width
- [ ] FS ! n ~ Set print mode(s) for Kanji characters
- [ ] FS & ~ Select Kanji character mode
- [ ] FS - n ~ Turn underline mode on/off for Kanji characters
- [ ] FS . ~ Cancel Kanji character mode
- [ ] FS 2 c1 c2 d1...dk ~ Define user-defined Kanji characters
- [ ] FS S n1 n2 ~ Set left- and right-side Kanji character spacing
- [ ] FS W n ~ Turn quadruple-size mode on/off for Kanji characters

Undocumented?:

- [ ] GS P ~ Specify horizontal and vertical units
- [ ] GS A ~ auto status back

Not Implemented:
- [x] GS V m n ~ Select cut mode and cut paper
  - In testing it didn't do anything
