package queries

import (
	"log"
	"strconv"

	"sort"
	"sync"

	
)

// Define the struct to hold the totalAmount and count
type UserBookingData struct {
	TotalAmount      int
	Count            int
	BookingCounts    []MonthlyBookingCount // Changed to slice of MonthlyBookingCount
	TopTheaters      []TheaterCount
}

// Define the TheaterCount structure
type TheaterCount struct {
	Theater string
	Count   int
}

// Define the MonthlyBookingCount structure
type MonthlyBookingCount struct {
	Month string
	Count  int
}
type MovieCount struct{
	Movie string
	Count   int
}
type AdminBokingData struct{
	
	TotalBooking int
	TotalAmount  int
	BookingCounts    []MonthlyBookingCount // Changed to slice of MonthlyBookingCount
	TopTheaters      []TheaterCount
	TopMovies []MovieCount
}
type OwnerBookingData struct{
	TotalBooking int
	TotalAmount  int
	BookingCounts    []MonthlyBookingCount // Changed to slice of MonthlyBookingCount
	TopTheaters      []TheaterCount
	TopMovies []MovieCount
}
func GetAllUserBookingsData(user string) (UserBookingData, error) {
	// Fetch all user bookings
	bookingData, err := GetAllUserBookings(user)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return UserBookingData{}, err
	}

	// Create channels for concurrent processing
	totalAmountChannel := make(chan int)
	countChannel := make(chan int)
	bookingCountsChannel := make(chan map[string]int)
	theaterVisitCountsChannel := make(chan map[string]int)

	var wg sync.WaitGroup

	// Goroutine for total amount and count
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalAmount := 0
		count := 0
		for _, item := range bookingData {
			count++ // Increment the count variable
			price, err := strconv.Atoi(item.Price)
			if err != nil {
				log.Printf("Failed to convert price: %v", err)
				continue
			}
			// Add the price to the total amount
			totalAmount += price
		}
		totalAmountChannel <- totalAmount // Send total amount to channel
		countChannel <- count               // Send count to channel
	}()

	// Goroutine for booking counts per month
	wg.Add(1)
	go func() {
		defer wg.Done()
		bookingCounts := make(map[string]int)
		for _, item := range bookingData {
			// Use the CreatedAt field directly
			createdAt := item.CreatedAt
			// Format the date as Year-Month (e.g., "2024-09")
			month := createdAt.Format("2006-01")
			bookingCounts[month]++
		}
		bookingCountsChannel <- bookingCounts // Send booking counts to channel
	}()

	// Goroutine for theater visit counts
	wg.Add(1)
	go func() {
		defer wg.Done()
		theaterVisitCounts := make(map[string]int)
		for _, item := range bookingData {
			theaterVisitCounts[item.Theater]++
		}
		theaterVisitCountsChannel <- theaterVisitCounts // Send theater visit counts to channel
	}()

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(totalAmountChannel)
		close(countChannel)
		close(bookingCountsChannel)
		close(theaterVisitCountsChannel)
	}()

	// Collect results
	totalAmount := <-totalAmountChannel
	count := <-countChannel
	bookingCounts := <-bookingCountsChannel
	theaterVisitCounts := <-theaterVisitCountsChannel

	// Convert bookingCounts map to slice of MonthlyBookingCount
	var monthlyCounts []MonthlyBookingCount
	for month, count := range bookingCounts {
		monthlyCounts = append(monthlyCounts, MonthlyBookingCount{Month: month, Count: count})
	}
	sort.Slice(monthlyCounts, func(i, j int) bool {
		return monthlyCounts[i].Month < monthlyCounts[j].Month
	})

	// Sort theaters by visit count
	var sortedTheaters []TheaterCount
	for theater, count := range theaterVisitCounts {
		sortedTheaters = append(sortedTheaters, TheaterCount{Theater: theater, Count: count})
	}
	sort.Slice(sortedTheaters, func(i, j int) bool {
		return sortedTheaters[i].Count > sortedTheaters[j].Count
	})

	// Get top 3 theaters
	var topTheaters []TheaterCount
	for i, theater := range sortedTheaters {
		if i >= 3 {
			break
		}
		topTheaters = append(topTheaters, theater)
	}

	// Create and return the result as UserBookingData
	userBookingData := UserBookingData{
		TotalAmount:   totalAmount,
		Count:         count,
		BookingCounts: monthlyCounts,
		TopTheaters:   topTheaters,
	}

	return userBookingData, nil
}


func GetAdminBookingData() (AdminBokingData, error) {
	// Fetch all user bookings
	bookings, err := GetAllBookingsAdmin()
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return AdminBokingData{}, err // Return an empty struct and the error if fetching fails
	}

	// Create channels for concurrent processing
	priceChannel := make(chan int)
	bookingCountsChannel := make(chan map[string]int)
	theaterVisitCountsChannel := make(chan map[string]int)
	movieSeenCountsChannel := make(chan map[string]int)

	var wg sync.WaitGroup

	// Goroutine for total price
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalPrice := 0
		for _, item := range bookings {
			price, err := strconv.Atoi(item.Price)
			if err != nil {
				log.Printf("Failed to convert price: %v", err)
				continue
			}
			totalPrice += price
		}
		priceChannel <- totalPrice // Send total price to channel
	}()

	// Goroutine for booking counts per month
	wg.Add(1)
	go func() {
		defer wg.Done()
		bookingCounts := make(map[string]int)
		for _, item := range bookings {
			// Use the CreatedAt field directly
			createdAt := item.CreatedAt
			// Format the date as Year-Month (e.g., "2024-09")
			month := createdAt.Format("2006-01")
			bookingCounts[month]++
		}
		bookingCountsChannel <- bookingCounts // Send booking counts to channel
	}()

	// Goroutine for theater visit counts
	wg.Add(1)
	go func() {
		defer wg.Done()
		theaterVisitCounts := make(map[string]int)
		for _, item := range bookings {
			theaterVisitCounts[item.Theater]++
		}
		theaterVisitCountsChannel <- theaterVisitCounts // Send theater visit counts to channel
	}()

	// Goroutine for movie seen counts
	wg.Add(1)
	go func() {
		defer wg.Done()
		movieSeenCounts := make(map[string]int)
		for _, item := range bookings {
			movieSeenCounts[item.Movie]++
		}
		movieSeenCountsChannel <- movieSeenCounts // Send movie seen counts to channel
	}()

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(priceChannel)
		close(bookingCountsChannel)
		close(theaterVisitCountsChannel)
		close(movieSeenCountsChannel)
	}()

	// Collect results
	totalPrice := <-priceChannel
	bookingCounts := <-bookingCountsChannel
	theaterVisitCounts := <-theaterVisitCountsChannel
	movieSeenCounts := <-movieSeenCountsChannel

	// Convert bookingCounts map to slice of MonthlyBookingCount
	var monthlyCounts []MonthlyBookingCount
	for month, count := range bookingCounts {
		monthlyCounts = append(monthlyCounts, MonthlyBookingCount{Month: month, Count: count})
	}
	sort.Slice(monthlyCounts, func(i, j int) bool {
		return monthlyCounts[i].Month < monthlyCounts[j].Month
	})

	// Sort theaters by visit count
	var sortedTheaters []TheaterCount
	for theater, count := range theaterVisitCounts {
		sortedTheaters = append(sortedTheaters, TheaterCount{Theater: theater, Count: count})
	}
	sort.Slice(sortedTheaters, func(i, j int) bool {
		return sortedTheaters[i].Count > sortedTheaters[j].Count
	})

	// Sort movies by seen count
	var sortedMovies []MovieCount
	for movie, count := range movieSeenCounts {
		sortedMovies = append(sortedMovies, MovieCount{Movie: movie, Count: count})
	}
	sort.Slice(sortedMovies, func(i, j int) bool {
		return sortedMovies[i].Count > sortedMovies[j].Count
	})

	// Get top 3 theaters
	var topTheaters []TheaterCount
	for i, theater := range sortedTheaters {
		if i >= 3 {
			break
		}
		topTheaters = append(topTheaters, theater)
	}

	// Get top 3 movies
	var topMovies []MovieCount
	for i, movie := range sortedMovies {
		if i >= 3 {
			break
		}
		topMovies = append(topMovies, movie)
	}

	// Create an AdminBookingData struct with the fetched bookings
	adminBookingData := AdminBokingData{
		TotalBooking: len(bookings),
		TotalAmount:  totalPrice,
		BookingCounts: monthlyCounts,
		TopTheaters:   topTheaters,
		TopMovies:     topMovies,
	}

	return adminBookingData, nil // Return the struct and no error if successful
}


func GetAllOwnerData(theaterNames []string) (OwnerBookingData, error) {
	// Fetch all user bookings
	bookings, err := GetOwnerReport(theaterNames)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return OwnerBookingData{}, err // Return an empty struct and the error if fetching fails
	}

	// Create channels for concurrent processing
	priceChannel := make(chan int)
	bookingCountsChannel := make(chan map[string]int)
	theaterVisitCountsChannel := make(chan map[string]int)
	movieSeenCountsChannel := make(chan map[string]int)

	var wg sync.WaitGroup

	// Goroutine for total price
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalPrice := 0
		for _, item := range bookings {
			price, err := strconv.Atoi(item.Price)
			if err != nil {
				log.Printf("Failed to convert price: %v", err)
				continue
			}
			totalPrice += price
		}
		priceChannel <- totalPrice // Send total price to channel
	}()

	// Goroutine for booking counts per month
	wg.Add(1)
	go func() {
		defer wg.Done()
		bookingCounts := make(map[string]int)
		for _, item := range bookings {
			// Use the CreatedAt field directly
			createdAt := item.CreatedAt
			// Format the date as Year-Month (e.g., "2024-09")
			month := createdAt.Format("2006-01")
			bookingCounts[month]++
		}
		bookingCountsChannel <- bookingCounts // Send booking counts to channel
	}()

	// Goroutine for theater visit counts
	wg.Add(1)
	go func() {
		defer wg.Done()
		theaterVisitCounts := make(map[string]int)
		for _, item := range bookings {
			theaterVisitCounts[item.Theater]++
		}
		theaterVisitCountsChannel <- theaterVisitCounts // Send theater visit counts to channel
	}()

	// Goroutine for movie seen counts
	wg.Add(1)
	go func() {
		defer wg.Done()
		movieSeenCounts := make(map[string]int)
		for _, item := range bookings {
			movieSeenCounts[item.Movie]++
		}
		movieSeenCountsChannel <- movieSeenCounts // Send movie seen counts to channel
	}()

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(priceChannel)
		close(bookingCountsChannel)
		close(theaterVisitCountsChannel)
		close(movieSeenCountsChannel)
	}()

	// Collect results
	totalPrice := <-priceChannel
	bookingCounts := <-bookingCountsChannel
	theaterVisitCounts := <-theaterVisitCountsChannel
	movieSeenCounts := <-movieSeenCountsChannel

	// Convert bookingCounts map to slice of MonthlyBookingCount
	var monthlyCounts []MonthlyBookingCount
	for month, count := range bookingCounts {
		monthlyCounts = append(monthlyCounts, MonthlyBookingCount{Month: month, Count: count})
	}
	sort.Slice(monthlyCounts, func(i, j int) bool {
		return monthlyCounts[i].Month < monthlyCounts[j].Month
	})

	// Sort theaters by visit count
	var sortedTheaters []TheaterCount
	for theater, count := range theaterVisitCounts {
		sortedTheaters = append(sortedTheaters, TheaterCount{Theater: theater, Count: count})
	}
	sort.Slice(sortedTheaters, func(i, j int) bool {
		return sortedTheaters[i].Count > sortedTheaters[j].Count
	})

	// Sort movies by seen count
	var sortedMovies []MovieCount
	for movie, count := range movieSeenCounts {
		sortedMovies = append(sortedMovies, MovieCount{Movie: movie, Count: count})
	}
	sort.Slice(sortedMovies, func(i, j int) bool {
		return sortedMovies[i].Count > sortedMovies[j].Count
	})

	// Get top 3 theaters
	var topTheaters []TheaterCount
	for i, theater := range sortedTheaters {
		if i >= 3 {
			break
		}
		topTheaters = append(topTheaters, theater)
	}

	// Get top 3 movies
	var topMovies []MovieCount
	for i, movie := range sortedMovies {
		if i >= 3 {
			break
		}
		topMovies = append(topMovies, movie)
	}

	// Create an OwnerBookingData struct with the fetched bookings
	ownerBookingData := OwnerBookingData{
		TotalBooking: len(bookings),
		TotalAmount:  totalPrice,
		BookingCounts: monthlyCounts,
		TopTheaters:   topTheaters,
		TopMovies:     topMovies,
	}

	return ownerBookingData, nil // Return the struct and no error if successful
}
