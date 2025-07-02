# stash

Stash is a lightweight file manager and provides session-based access control.

## Features

- File upload with progress tracking and speed estimation
- File listing with metadata (size, modified time)
- Download and delete functionality
- Search, sorting, and folder navigation

## Routes
### UI

* `GET /` — File manager (requires authentication)
* `GET /files` — Partial render of file list
* `POST /upload` — File upload
* `GET /download/:filename` — File download
* `DELETE /delete/:filename` — File deletion

### Static

* `/assets` — Static assets (CSS, JS, icons)

Authentication routes are managed by the external `snare` module.

## Security Notes

* Uploaded files are stored under `/home/public`
* Only authenticated users can access file listing and actions
* Cookies are HTTP-only and session-based
* File names are not sanitized;

## License

MIT License
