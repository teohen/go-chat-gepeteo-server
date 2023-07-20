
# Change Log
All notable changes to this project will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [1.3.0] - 2023-07-19

### Added
 - Added the server ls_clients ROOM_NAME command
 - Added the client join ROOM_NAME command
 - Added the requirement to not allow creation of 2 rooms with the same name
 - Added the requirement to not allow registration of 2 clients with the same name


### Changed
 - Refactored the server structure moving parts to their module (server, room, client, etc)

### Fixed
- Fixed the bug where every new client connection creates a new general room

