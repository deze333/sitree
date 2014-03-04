package sitree

import (
    "fmt"
    "time"
)

//------------------------------------------------------------
// Scheduler Model
//------------------------------------------------------------

type Scheduler struct {
	day, hour        int
	timer            *time.Timer
	executor         func()
}

//------------------------------------------------------------
// Scheduler Methods
//------------------------------------------------------------

// Schedules repeated call to server processing function.
// day -- step in days, >= 0.
// hour -- specific hour of the day (if day > 0) or every Nth hour (if day == 0).
func (s *Scheduler) Set(day, hour int, fn func()) (err error) {
    if day < 0 {
        return fmt.Errorf("[sitree] Scheduler: Cannot use negative day amount: %v", day)
    }
    if day == 0 {
        if hour < 1 || hour > 24 {
            return fmt.Errorf("[sitree] Scheduler: Hour must be within [1, 24] range: %v", hour)
        }
    } else {
        if hour < 0 || hour > 24 {
            return fmt.Errorf("[sitree] Scheduler: Hour must be within [0, 23] range: %v", hour)
        }
    }
    s.day = day
    s.hour = hour
    s.executor = fn
    s.timer = time.AfterFunc(timerDelta(time.Now(), s.day, s.hour), s.periodic)

    return nil
}

// Executor function that runs when timer is up.
func (s *Scheduler) periodic() {
    fmt.Println("[sitree]  *** EXECUTING SCHEDULER PERIODIC ***")
    s.executor()
    s.timer = time.AfterFunc(timerDelta(time.Now(), s.day, s.hour), s.periodic)
}

// Calculates delta time between 'now' time and scheduled recurrent time
// speicified by day and hour.
func timerDelta(now time.Time, day, hour int) time.Duration {
    // Run every Nth hour, day not taken into account
    if day == 0 {
        // Force correct negative hour
        if hour <= 0 {
            hour = 1
        }
        //t2 := now.Add(time.Minute * time.Duration(hour)) // DEBUG! 
        t2 := now.Add(time.Hour * time.Duration(hour))
        t2 = time.Date(
            t2.Year(),
            t2.Month(),
            t2.Day(),
            //t2.Hour(), t2.Minute(), 0, 0, // DEBUG!
            t2.Hour(), 00, 0, 0,
            t2.Location())

        fmt.Println("[sitree] Next run:", t2)
        return t2.Sub(now)
    }

    // Run once every Nth day at specific hour.
    // Force correct negative days
    if day < 0 {
        day = 1
    }

    // Fast forward N days
    t2 := now.AddDate(0, 0, day)
    // Set the hour for that future day
    t2 = time.Date(
        t2.Year(),
        t2.Month(),
        t2.Day(),
        hour, 00, 0, 0,
        t2.Location())

    fmt.Println("[sitree] Next run:", t2)
    return t2.Sub(now)
}
