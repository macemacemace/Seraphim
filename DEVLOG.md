# Seraphim development log

A short log of what I worked on, what I decided, where I got stuck and where AI helped or got things wrong.

**How this log is written:** I write all the app code myself, with Claude (Claude Code) as a tutor that explains, reviews my code and checks it. Claude keeps notes after every session, and drafts the entries in this log from those notes every few sessions. I read every entry before committing it.

Entries up to 15 Nov 2026 were written on 15 Nov from those session notes, when I received the official project guidelines and learned a log was required. From then on the log is updated as the project goes.

---

## 28 Sep - 5 Oct: setup

- Set up the monorepo: Go backend in `backend/`, React + TypeScript (Vite) in `frontend/`, Postgres with PostGIS in Docker, GitHub Actions CI.
- **Decision: a real Go backend instead of a mock.** Only the frontend is graded, but I wanted to learn Go and WebSockets, and a real API makes the frontend work against real data.
- **Decision: React instead of Angular.** I already used React in my Eclipse project and wanted to go deeper with it rather than learn two new things at once (Go and Angular).
- **Decision: tooling.** goose for migrations, sqlc to generate typed Go code from SQL (instead of an ORM, so I write and understand the SQL myself).
- First week of the plan was finished on day 1.

## 29 Sep - 17 Oct: authentication

- Register, login with a JWT access token, auth middleware, `GET /api/me`, register and login pages with styling.
- **Decision: build auth by hand (bcrypt + JWT) instead of Supabase Auth.** I used Supabase Auth before, and building it myself was the point.
- **Decision: any email can register.** Many students don't have a university email, and graduates selling books are welcome too.
- **Decision: refresh tokens and email verification moved to later,** so listings could start on time.
- Chose the colour palette and set it up as CSS variables from the start, mobile-first.
- Stuck on: CI went red only because a file was not formatted with gofmt. I learned to format before committing.
- Stuck on: I worked on a feature branch on my laptop and my PC didn't have it, because sync only pulls the current branch. Decided to work on main.
- My mistakes: VS Code auto-imported a crypto package (`secretbox`) when I typed `secret`; `=` instead of `==`; I copied the register page for login and forgot to change the URL.
- AI helped: explaining `errors.As` (compared it to `instanceof` in JavaScript) and Go error wrapping.

## 18 - 27 Oct: listings, feed, sell and detail pages

- Listings table with a PostGIS location, endpoints to create a listing, get one, and a feed of listings within 50 km.
- **Decision: seller location is rounded to about 1 km before it is stored,** so the exact home address is never saved.
- **Decision: if the user blocks location, the feed shows all listings, newest first,** instead of an error page.
- Feed page with the browser geolocation API, a card grid, loading dots and an empty state. Sell page and listing detail page (`/listings/:id`).
- **Decision: nav bar postponed** until messaging, because messages need a badge in it. The feed got its own Sell button for now.
- Stuck on: location was "automatically blocked" in Chrome, so the feed never loaded nearby listings. Resetting the site permission fixed it.
- Stuck on: I pushed a broken commit to main once (CI red) and fixed it the same evening.
- My mistakes: one condition in the Sell page set the wrong variable, and `eslint --fix` turned it into a `const` and hid the bug. Lesson: auto-fix silences warnings, it doesn't fix logic.
- My mistakes: `==` instead of `===` in TypeScript, code typed into the wrong file.

## 28 Oct - 5 Nov: photos

- **Decision: Cloudinary for photo storage.** MinIO no longer has a free tier. I compared Cloudinary, Cloudflare R2 and Supabase Storage and picked Cloudinary.
- Up to 5 photos per listing, required for items, optional for services. Max 5 MB, JPG, PNG or WebP only. The Go API checks the real file type from the first bytes of the file, not just the file name.
- Upload endpoint, file picker on the Sell page, photos on the detail page with arrows and thumbnails.
- **Decision: the arrows wrap around** (next on the last photo goes to the first).
- Thumbnails are buttons, not plain images, so they work with the keyboard.
- Stuck on: my ownership check compared the listing ID instead of the seller ID, so anyone could add photos to any listing. Caught it in review and fixed it before it was used.
- Stuck on: new backend code "didn't work" because an old Go server was still running on the port and answering instead of the new one. I learned to check which process holds the port.
- AI got it wrong: Claude told me my photo viewer code was deleted, because it only read part of the file. It wasn't.
- Took 30 Oct - 4 Nov off for midterms.

## 5 - 8 Nov: cover photos and search

- First photo of each listing on the feed cards, using one SQL subquery instead of one extra request per card (avoiding the N+1 problem).
- Search in title and description.
- **Decision: search lives in the URL (`/?q=laptop`)**, so a search can be shared and the back button works.
- "Add photos" button instead of the default file input. The real input is visually hidden but still keyboard accessible.
- Tested the bad cases of photo upload: no token (401), a .txt file (400), someone else's listing (403), missing listing (404).
- Stuck on: sqlc said a column was "ambiguous" after I added the subquery. Fixed by naming the table (`listings.created_at`).
- AI got it wrong: Claude first told me to use `hidden` on the file input, which would have made it unreachable with the keyboard. Corrected to `opacity: 0`.
- AI got it wrong: Claude gave me `git add` paths relative to the wrong folder.

## 8 - 13 Nov: filters, sort, hot reload, profile

- Filters by type (item or service) and price range, and sort (newest, oldest, price up, price down). All of them live in the URL like search.
- **Decision: max price uses -1 for "no limit", not 0,** because a max price of 0 is a real search (free items). The API rejects negative prices from the URL, so -1 can't be sent by a user.
- Set up air for hot reload of the Go server, so the old-server-on-the-port problem is gone.
- Avatar button with the first letter of the name, profile page with log out, and an endpoint for my own listings.
- **Decision: the profile has no photo upload,** just a letter avatar. Google sign-in comes after messaging.
- My mistakes: I cut a SQL block instead of copying it, so one of the two feed queries lost the filters; a typo (`filers`) got copied to 5 lines; I wrote the max price check with the sign the wrong way round twice.
- AI got it wrong: Claude said both feed queries had the filter block when only one did, because it only read the top half of the file. Since then it reads whole files when checking.
- Stuck on: after a sync script ran, Vite kept serving an old version of the feed page, so the filters looked ignored even though my code was right. Saving the file again fixed it.

## 15 Nov: project guidelines

- Received the official guidelines on 13 Nov and compared the project against them.
- Missing: this log, a proposal (27 Nov), a deployed URL (14 Dec checkpoint), frontend tests, reusable components, error states for failed requests, README with AI disclosure, removing leftover starter files.
- **Decision: state management is `useState` plus the URL.** Search, filters and sort are URL search params (`useSearchParams`), so links can be shared and the back button works. Everything else is local component state. I didn't add Redux or a global store because no state is shared between pages yet. When the nav bar needs the logged-in user on every page, React Context is the likely next step.
- **Plan:** tests and README at the end before the deadline; reusable components and error states while finishing the profile page.
