package bytestostr

import "testing"

var text = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
	"Vivamus accumsan mauris sit amet tempor lacinia. Nam elementum dictum pharetra. " +
	"Donec arcu turpis, aliquet ut risus in, pulvinar ornare sem. Suspendisse " +
	"lobortis massa sit amet velit hendrerit, id ullamcorper est semper. " +
	"In sodales nibh a laoreet blandit. Nullam sit amet faucibus dui. " +
	"Donec faucibus ante sed egestas tempor. Phasellus luctus aliquet " +
	"tincidunt. Pellentesque faucibus dignissim tempus. Quisque eu " +
	"euismod ligula.\n\nSed ultrices lorem ut pellentesque sodales. " +
	"Nunc luctus turpis a tellus tempor, cursus lobortis sapien rutrum. " +
	"Cras dolor nisl, luctus ac tempor vitae, mollis eu diam. " +
	"Pellentesque commodo laoreet cursus. Cras non nisl leo. " +
	"Donec hendrerit nulla lectus, ac congue velit semper nec. " +
	"Duis sit amet odio ligula. Pellentesque consequat est quis " +
	"magna tempor pharetra.\n\nFusce at ultricies purus. Integer " +
	"in enim quam. In congue nibh ante, ut porttitor ante " +
	"volutpat ac. Fusce pellentesque semper augue imperdiet porta. " +
	"Ut scelerisque imperdiet arcu et eleifend. Nam sit amet dui " +
	"consectetur, sodales metus non, porta arcu. Integer et elit " +
	"in lacus finibus semper. Aenean egestas sapien nec rutrum dapibus. " +
	"In vel est velit. Sed sed feugiat justo, pretium interdum neque. " +
	"Sed velit nisl, tincidunt sit amet erat sed, finibus consequat " +
	"libero. Vivamus pulvinar id nisi eu iaculis.\n\n")

func BenchmarkSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = string(text)
	}
}

func BenchmarkBytesToStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BytesToStr(text)
	}
}
