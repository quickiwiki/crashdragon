# Changelog

## Version x.x.x
* Add: API endpoints
* Add: `-version` command line argument, which prints the current CrashDragon version, then exits
* Add: `Housekeeping.ReportRetentionTime` config field to configure the duration Reports should be kept (supports either integers or Go duration strings). Reports older then this retention time will be removed from the database and disk.
* Fix: Version selector now stays in the right order after deleting a version, #48
* Fix: Move Crash counts to own table to get rid of numerous `count()`s on the main page, see #48
* Change: Use Go Modules instead of third party govendor
* Change: Module name trimming is now configurable (on or off) with the `Symbolicator.TrimModuleNames` config option
* Change: `ProcessUptime` now uses `uint64`, see !7
* Change: Update `breakdpad`-submodule to current `master`
* Remove: `-config` command line argument in favor of dedicated config-search-paths (`/etc/crashdragon/config.ext`, `$HOME/.crashdragon/config.ext`, `./config.ext`)

## Version 1.2.1

* Fix: Commenting on crashes possible again
* Fix: Module should match now correctly to signatures

## Version 1.2.0

* Add: `-config`-flag to specify config file when running CrashDragon
* Add: Module database field for Reports and Crashes (See #39)
* Change: Now gracefully shutting down on `SIGINT` (See #25)
* Change: Only show reports matching the filtered version when in Crash-view
* Change: Allow changing Product and Version filters in detail views
* Change: Fixed now with date instead of boolean (See #36)
* Change: Symbol files now in seperate file locations based on product and version (See #15)
* Fix: Show correct crash counts when filtered by version (See #40)
* Fix: Installation of stylesheets

## Version 1.1.1

* Add: `govendor sync` step to update guide
* Change: Updated GORM to current master
* Fix: Crash-Version-Association
* Fix: Correctly install stylesheet

## Version 1.1.0

* Add: Version dropdown next to product dropdown
* Add: Button to delete reports and associated minidump and text files
* Add: Symbolicator path configurable in the config file
* Add: Install target to Makefile **(Read updated UPDATE guide!)**
* Add: Changelog
* Change: Show more details of a crahs on the crash page
* Change: Updated UPDATE guide
* Change: Updated INSTALL guide
* Fix: Create version if symbols are uploaded for non-existing version
* Fix: Crash on wrong symbolicator output (CrashingThread > len(Threads))

## Version 1.0.3

* Add: Pagination on crashes page
* Change: Run migrations as transaction
* Change: Improve UPDATE guide
* Change: Update breakpad submodule
* Change: Order crashes by popularity
* Change: Store generated TXT reports in file system
* Change: Only return ID to report sender (instead of whole json object)
* Fix: Missing signature in report view
* Fix: Do not unnecessarily generate TXT reports
* Fix: Fail after first DB connection try
* Fix: Performance improvements
* Fix: Show full crashed thread instead of truncated version
* Fix: Check for errors on `io.Copy()`
* Fix: Migrations naming
* Remove: Obsolete upload_syms application

## Version 1.0.2

* Fix: Missing govendor dependency

## Version 1.0.1

* Add: UPDATE guide
* Add: Migrations
* Change: Move API endpoint `/api/` to `/api/v1/`
* Change: Serve Chart.js locally

## Version 1.0.0
Initial release
