package mainwindow

import (
	"log"

	"github.com/funcan/soapyradiotool/mathtools"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

func Setup(registerListener func(string, func(interface {})), registerSource func(string) func(interface {})) {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Analyser")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	pagegrid, err := gtk.GridNew()
	if err != nil {
                log.Fatal("Unable to create page grid:", err)
        }
	pagegrid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	l, err := gtk.LabelNew("Filename: ")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	loadFileSender := registerSource("load file")
	processingSteps := make([]func([]float64)[]float64, 0)

	filenamebox, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create filenamebox:", err)
	}
	filenamebox.Connect("activate", func(e interface{}) {
		filename, _ := filenamebox.GetText()
		loadFileSender(filename)
	})

	loadbutton, err := gtk.ButtonNewWithLabel("Load")
	if err != nil {
		log.Fatal("Unable to create filenamebox:", err)
	}
	loadbutton.Connect("clicked", func() {
		filename, _ := filenamebox.GetText()
		processingSteps = make([]func([]float64)[]float64, 0)
		loadFileSender(filename)
	})

	loadgrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	loadgrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	loadgrid.Add(l)
	loadgrid.Add(filenamebox)
	loadgrid.Add(loadbutton)

	pagegrid.Add(loadgrid)

	chartarea, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("Unable to create drawing area:", err)
	}

	var dataPtr *[]float64

	registerListener("show data", func(e interface{}){
		data, ok := e.([]float64)
		if !ok {
			log.Printf("show data for chartarea bad param %T", e)
			return
		}
		dataPtr = &data
		chartarea.QueueDraw()
	})

	pagegrid.Add(chartarea)
	chartarea.SetHExpand(true)
	chartarea.SetVExpand(true)

	adjustment, err := gtk.AdjustmentNew(0, 0, 1, 1, 1, 1)
	if err != nil {
		log.Fatal("Unable to create adjustment:", err)
	}

	adjustment.Connect("value-changed", func() {
		chartarea.QueueDraw()
	})

	registerListener("show data", func(e interface{}){
		data, ok := e.([]float64)
                if !ok {
                        log.Printf("show data for adjustment bad param %T", e)
                        return
                }
		adjustment.Configure(0, 0, float64(len(data)), 10, 1, 100)
	})

	scrollbar, err := gtk.ScrollbarNew(gtk.ORIENTATION_HORIZONTAL, adjustment)
	if err != nil {
		log.Fatal("Unable to create scrollbar:", err)
	}
	scrollbar.SetHExpand(true)
	pagegrid.Add(scrollbar)

	zoom := 1

	chartarea.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		if dataPtr != nil {
			drawHandler(da, cr, *dataPtr, adjustment, zoom, processingSteps)
		}
	})


	controlsgrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create controls grid:", err)
	}
	controlsgrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)

	zoombutton, err := gtk.SpinButtonNewWithRange(1, 10, 1);
	if err != nil {
		log.Fatal("Unable to create zoom button:", err)
	}
	zoombutton.Connect("value-changed", func() {
		zoom = zoombutton.GetValueAsInt()
		chartarea.QueueDraw()
	})

	controlsgrid.Add(zoombutton)

	absButton, err := gtk.ButtonNewWithLabel("Abs Around Mean")
	if err != nil {
                log.Fatal("Unable to create abs button")
	}

	absButton.Connect("clicked", func() {
		processingSteps = append(processingSteps, func(in []float64)[]float64 {
			return mathtools.AbsAroundMean(in)
		})
		chartarea.QueueDraw()
	})
	controlsgrid.Add(absButton)

	// FIXME: Make floor adjustable
	squelchButton, err := gtk.ButtonNewWithLabel("Squelch")
	if err != nil {
		log.Fatal("Unable to create squelch button")
	}
	squelchButton.Connect("clicked", func() {
                processingSteps = append(processingSteps, func(in []float64)[]float64 {
                        return mathtools.Squelch(in, mathtools.StdDev(in)*2)
		})
		chartarea.QueueDraw()
	})
	controlsgrid.Add(squelchButton)

	denoiseButton, err := gtk.ButtonNewWithLabel("Denoise")
	if err != nil {
                log.Fatal("Unable to create denoise button")
	}
	denoiseButton.Connect("clicked", func() {
                processingSteps = append(processingSteps, func(in []float64)[]float64 {
                        return mathtools.Denoise(in)
                })
                chartarea.QueueDraw()
        })
	controlsgrid.Add(denoiseButton)

	// FIXME: Make buckets adjustable
	rollingAvgButton, err := gtk.ButtonNewWithLabel("Rolling Average")
	if err != nil {
                log.Fatal("Unable to create rolling avg button")
	}

	rollingAvgButton.Connect("clicked", func() {
                processingSteps = append(processingSteps, func(in []float64)[]float64 {
                        return mathtools.RollingAverage(in, 5)
                })
                chartarea.QueueDraw()
        })
	controlsgrid.Add(rollingAvgButton)

	// FIXME: Make buckets adjustable
	edgeFinderButton, err := gtk.ButtonNewWithLabel("Edge finder")
	if err != nil {
                log.Fatal("Unable to create edge finder button")
	}

	edgeFinderButton.Connect("clicked", func() {
                processingSteps = append(processingSteps, func(in []float64)[]float64 {
                        return mathtools.EdgeFinder(in, 5)
                })
                chartarea.QueueDraw()
        })
	controlsgrid.Add(edgeFinderButton)

	pagegrid.Add(controlsgrid)

	stepsgrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create sets grid:", err)
	}
	stepsgrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)

	dummybutton, err := gtk.ButtonNewWithLabel("")
	if err != nil {
		log.Fatal("Unable to create dummybutton:", err)
	}
	stepsgrid.Add(dummybutton)

	pagegrid.Add(stepsgrid)

	win.Add(pagegrid)

	win.SetDefaultSize(800, 600)
	win.ShowAll()

	gtk.Main()
}

func drawHandler(da *gtk.DrawingArea, cr *cairo.Context, data []float64, adjustment *gtk.Adjustment, zoom int, processingsteps []func([]float64)[]float64) {
	width := da.GetAllocatedWidth()
	height := da.GetAllocatedHeight()

	start := int(adjustment.GetValue())
	if start+width > len(data) {
		start = len(data) - width
	}
	if start < 0 {
		start = 0
	}

	for _, processingstep := range(processingsteps) {
		data = processingstep(data)
	}

	max := mathtools.Max(data)

	cr.SetSourceRGBA(1,0,0,1)
	cr.SetLineWidth(0.6)
	if zoom == 1 {
		// 3 pixels per sample
		for i:=0; i<width/3; i++ {
			sample := data[i+start]
			scaled := float64(height) - ((sample/max) * float64(height))
			if i==0 {
				cr.MoveTo(float64(i*3), scaled)
			} else {
				cr.LineTo(float64(i*3), scaled)
			}
		}
	} else {
		zoom2steps := map[int]int {
			2: 1,
			3: 2,
			4: 5,
			5: 10,
			6: 25,
			7: 75,
			8: 200,
			9: 500,
			10: 1000,
		}
		steps := zoom2steps[zoom]

		for i:=0; i<width; i++ {
			// FIXME: Fix nasty sampling/aliasing problem
			value := (i*steps) + start
			if value >= len(data) {
				continue
			}
			sample := data[value]
			scaled := float64(height) - ((sample/max) * float64(height))
			if i==0 {
				cr.MoveTo(float64(i), scaled)
			} else {
				cr.LineTo(float64(i), scaled)
			}
		}
	}
	cr.Stroke()
}
