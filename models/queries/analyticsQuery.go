package queries

import (
	"log"
	"strconv"

	"sort"

	
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

	totalAmount := 0
	count := 0
	bookingCounts := make(map[string]int)
	theaterVisitCounts := make(map[string]int)

	for _, item := range bookingData {
		count++ // Increment the count variable
		price, err := strconv.Atoi(item.Price)
		if err != nil {
			log.Printf("Failed to convert price: %v", err)
			continue
		}
		// Add the price to the total amount
		totalAmount += price
		// Use the CreatedAt field directly
		createdAt := item.CreatedAt

		// Format the date as Year-Month (e.g., "2024-09")
		month := createdAt.Format("2006-01")

		// Increment the count for this month
		bookingCounts[month]++
		// Count the visits to each theater
		theaterVisitCounts[item.Theater]++
	}

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

func GetAdminBookingData()(AdminBokingData,error) {
	// Fetch all user bookings
	bookings, err := GetAllBookingsAdmin()
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return AdminBokingData{}, err // Return an empty struct and the error if fetching fails
	}
	
totalPrice:=0
bookingCounts := make(map[string]int)
	theaterVisitCounts := make(map[string]int)
	movieSeenCounts := make(map[string]int)
for _,item:=range bookings{
	price, err := strconv.Atoi(item.Price)
		if err != nil {
			log.Printf("Failed to convert price: %v", err)
			continue
		}
		totalPrice += price
		// Use the CreatedAt field directly
		createdAt := item.CreatedAt

		// Format the date as Year-Month (e.g., "2024-09")
		month := createdAt.Format("2006-01")

		// Increment the count for this month
		bookingCounts[month]++
		// Count the visits to each theater
		theaterVisitCounts[item.Theater]++
		movieSeenCounts[item.Movie]++

}
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
	// sort movies by visit count
	var sortedMovies []MovieCount
	for movie,count := range movieSeenCounts{
		sortedMovies = append(sortedMovies,MovieCount{Movie: movie,Count: count} )
	}
	sort.Slice(sortedTheaters, func(i, j int) bool {
		return sortedTheaters[i].Count > sortedTheaters[j].Count
	})
sort.Slice(sortedMovies,func(i, j int) bool {
	return sortedMovies[i].Count>sortedMovies[j].Count
})
	// Get top 3 theaters
	var topTheaters []TheaterCount
	for i, theater := range sortedTheaters {
		if i >= 3 {
			break
		}
		topTheaters = append(topTheaters, theater)
	}
	//get top 3 movies
	var topMovies []MovieCount
	for i,movie:= range sortedMovies{
		if i>=3{
			break
		}
		topMovies = append(topMovies,movie)
	}
	// Create an AdminBookingData struct with the fetched bookings
	adminBookingData := AdminBokingData{
		
		TotalBooking: len(bookings),
		TotalAmount: totalPrice,
		BookingCounts: monthlyCounts,
		TopTheaters:   topTheaters,
		TopMovies: topMovies,
	}



	return adminBookingData, nil // Return the struct and no error if successful
}


func GetAllOwnerData(theaterNames []string )(OwnerBookingData,error){
	// Fetch all user bookings
	bookings, err := GetOwnerReport(theaterNames)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return OwnerBookingData{}, err // Return an empty struct and the error if fetching fails
	}
	
totalPrice:=0
bookingCounts := make(map[string]int)
	theaterVisitCounts := make(map[string]int)
	movieSeenCounts := make(map[string]int)
for _,item:=range bookings{
	price, err := strconv.Atoi(item.Price)
		if err != nil {
			log.Printf("Failed to convert price: %v", err)
			continue
		}
		totalPrice += price
		// Use the CreatedAt field directly
		createdAt := item.CreatedAt

		// Format the date as Year-Month (e.g., "2024-09")
		month := createdAt.Format("2006-01")

		// Increment the count for this month
		bookingCounts[month]++
		// Count the visits to each theater
		theaterVisitCounts[item.Theater]++
		movieSeenCounts[item.Movie]++

}
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
	// sort movies by visit count
	var sortedMovies []MovieCount
	for movie,count := range movieSeenCounts{
		sortedMovies = append(sortedMovies,MovieCount{Movie: movie,Count: count} )
	}
	sort.Slice(sortedTheaters, func(i, j int) bool {
		return sortedTheaters[i].Count > sortedTheaters[j].Count
	})
sort.Slice(sortedMovies,func(i, j int) bool {
	return sortedMovies[i].Count>sortedMovies[j].Count
})
	// Get top 3 theaters
	var topTheaters []TheaterCount
	for i, theater := range sortedTheaters {
		if i >= 3 {
			break
		}
		topTheaters = append(topTheaters, theater)
	}
	//get top 3 movies
	var topMovies []MovieCount
	for i,movie:= range sortedMovies{
		if i>=3{
			break
		}
		topMovies = append(topMovies,movie)
	}
	// Create an AdminBookingData struct with the fetched bookings
	ownerBookingData := OwnerBookingData{
		
		TotalBooking: len(bookings),
		TotalAmount: totalPrice,
		BookingCounts: monthlyCounts,
		TopTheaters:   topTheaters,
		TopMovies: topMovies,
	}



	return ownerBookingData, nil // Return the struct and no error if successful
}