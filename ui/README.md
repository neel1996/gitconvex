# gitconvex react project
This is the front-end react source for [gitconvex](https://github.com/neel1996/gitconvex).

<p align="center">
    <img src="https://user-images.githubusercontent.com/47709856/87220396-e72df380-c380-11ea-9b2b-e156402842bb.png" width="280">
</p>

## Dependencies

The dependency packages used by this project can be found [here](https://github.com/neel1996/gitconvex-ui/network/dependencies)

- **🎨 Styling** - For styling, the project uses [tailwind](https://github.com/tailwindlabs/tailwindcss) css framework 
- **📑 Syntax Highlighting** - [prismjs](https://github.com/PrismJS/prism) is used for syntax highlighting within the *Git Difference* and *CodeView* section
- **🎭 Icon set** - [FontAweomse for react](https://github.com/FortAwesome/Font-Awesome)

## Contribute!

If you are interested in contributing to the project, fork the repo and submit a PR. Currently there are only 3 maintainers working on the project, so the review may take some time. Hopefully will get a couple more on board soon to maintain the repo

### Guidelines 

Fork the repo and raise a new Pull Request to merge your branch with the `development` branch of this repo. Once the review is complete, the PR will be approved and merged with `master`

##### API schema request

The backend is powered by graphql and if your changes require a separate query / mutation to be created to communicate with the server, then suggest
your schema in the [api_schema.graphql](api_schema.graphql) file.
### Setup

After cloning the repo, follow the steps mentioned below to setup the react app,

- **Installing dependencies**

``` shell
$ cd gitconvex-ui
$ npm install
```

- **Building the css file**

The project uses `tailwindcss v2.0.2` for styling all the elements, so it is mandatory to build the CSS file which is not included in the repo. Follow the steps to build the css file

``` shell

$ cd gitconvex-ui
$ npm install --global tailwindcss 

$ npm run build:tailwind

## This will generate a default tailwind css bundle

```

> **Note:** The final production build stage is configured to purge unused CSS selectors from the tailwind css file. So make sure you follow the [tailwind purge guidelines](https://tailwindcss.com/docs/controlling-file-size#writing-purgeable-html:~:text=Don't%20use%20string%20concatenation%20to%20create%20class%20names) to ensure that the required styles are included to the [production bundle](https://github.com/neel1996/gitconvex)

- **Starting the app**

After completing the setup process, use `npm start` to start the react app


## Project directory tree

``` shell
├── LICENSE
├── README.md
├── api_schema.graphql
├── package-lock.json
├── package.json
├── public
│   ├── favicon.ico
│   ├── gitconvex.png
│   ├── index.html
│   ├── logo192.png
│   ├── logo512.png
│   ├── manifest.json
│   ├── prism.css
│   └── robots.txt
├── src
│   ├── App.css
│   ├── App.js
│   ├── Components
│   │   ├── Animations
│   │   │   └── InfiniteLoader.js
│   │   ├── DashBoard
│   │   │   ├── AddNewRepoComponent
│   │   │   │   ├── AddNewRepoComponent.tsx
│   │   │   │   ├── AddNewRepoContainer.tsx
│   │   │   │   ├── AddRepoActionButtonsComponent.tsx
│   │   │   │   ├── AddRepoFormComponent.tsx
│   │   │   │   ├── AddRepoStatusAlert.tsx
│   │   │   │   ├── AuthComponent
│   │   │   │   │   ├── AuthComponent.tsx
│   │   │   │   │   ├── AuthOptionsComponent.tsx
│   │   │   │   │   ├── HTTPSAuthForm
│   │   │   │   │   │   ├── HTTPSAuthForm.tsx
│   │   │   │   │   │   └── HTTPSAuthHintComponent.tsx
│   │   │   │   │   └── SSHAuthForm.tsx
│   │   │   │   ├── CloneComponent.tsx
│   │   │   │   ├── ToggleSwitchComponent.tsx
│   │   │   │   └── add-new-repo-state
│   │   │   │       ├── actions.ts
│   │   │   │       ├── addRepoContext.tsx
│   │   │   │       └── addRepoReducer.ts
│   │   │   ├── Compare
│   │   │   │   ├── BranchCompareComponent
│   │   │   │   │   ├── BranchCommitLogChanges.js
│   │   │   │   │   └── BranchCompareComponent.js
│   │   │   │   ├── CommitCompareComponent
│   │   │   │   │   ├── CommitCompareComponent.js
│   │   │   │   │   ├── CommitFileDifferenceComponent.js
│   │   │   │   │   └── CommitLogCardComponent.js
│   │   │   │   ├── CompareActionButtons.js
│   │   │   │   ├── CompareActiveRepoPane.js
│   │   │   │   ├── CompareComponent.js
│   │   │   │   ├── CompareSelectionHint.js
│   │   │   │   ├── RepoSearchBar.js
│   │   │   │   └── SearchRepoCards.js
│   │   │   ├── Dashboard.js
│   │   │   ├── DashboardPaneComponents
│   │   │   │   ├── LeftPane.js
│   │   │   │   └── RightPane.js
│   │   │   ├── Help
│   │   │   │   └── Help.js
│   │   │   ├── Repository
│   │   │   │   ├── GitComponents
│   │   │   │   │   ├── GitDiffViewComponent.js
│   │   │   │   │   ├── GitOperation
│   │   │   │   │   │   ├── CommitComponent.js
│   │   │   │   │   │   ├── GitOperationComponent.js
│   │   │   │   │   │   ├── PushComponent.js
│   │   │   │   │   │   └── StageComponent.js
│   │   │   │   │   └── GitTrackedComponent.js
│   │   │   │   └── RepoComponents
│   │   │   │       ├── AddRepoForm.js
│   │   │   │       ├── RepoCard.js
│   │   │   │       ├── RepoComponent.js
│   │   │   │       ├── RepoDetails
│   │   │   │       │   ├── FileExplorerComponent.js
│   │   │   │       │   ├── RepoDetailBackdrop
│   │   │   │       │   │   ├── AddBranchComponent.js
│   │   │   │       │   │   ├── BranchListComponent.js
│   │   │   │       │   │   ├── CodeFileViewComponent.js
│   │   │   │       │   │   ├── CommitLogComponent.js
│   │   │   │       │   │   ├── CommitLogFileCard.js
│   │   │   │       │   │   ├── FetchPullActionComponent.js
│   │   │   │       │   │   ├── RemoteConfigComponent
│   │   │   │       │   │   │   ├── AddRemoteRepoFormComponent.js
│   │   │   │       │   │   │   ├── RemoteCard.js
│   │   │   │       │   │   │   └── RemoteManagementComponent.js
│   │   │   │       │   │   └── SwitchBranchComponent.js
│   │   │   │       │   ├── RepoInfoComponent.js
│   │   │   │       │   ├── RepoLeftPaneComponent.js
│   │   │   │       │   ├── RepoRightPaneComponent.js
│   │   │   │       │   ├── RepositoryDetails.js
│   │   │   │       │   └── backdropActionType.js
│   │   │   │       └── RepositoryAction.js
│   │   │   └── Settings
│   │   │       ├── Settings.js
│   │   │       ├── SettingsDataFileComponent.js
│   │   │       ├── SettingsPortComponent.js
│   │   │       └── SettingsRepoListComponent
│   │   │           ├── SettingsRepoDeleteBackDrop.js
│   │   │           ├── SettingsRepoDeleteComponent.js
│   │   │           ├── SettingsRepoListCard.js
│   │   │           └── SettingsRepoListComponent.js
│   │   ├── LoadingHOC.js
│   │   ├── SplashScreen.css
│   │   ├── SplashScreen.js
│   │   └── styles
│   │       └── LeftPane.css
│   ├── actionStore.js
│   ├── assets
│   │   ├── gitconvex.png
│   │   └── icons
│   ├── context.js
│   ├── index.css
│   ├── index.js
│   ├── prism.css
│   ├── react-app-env.d.ts
│   ├── reducer.js
│   ├── serviceWorker.js
│   ├── setupTests.js
│   ├── tailwind.config.js
│   ├── tests
│   │   ├── App.test.js
│   │   └── Dashboard.test.js
│   └── util
│       ├── apiURLSupplier.js
│       ├── env_config.js
│       └── relativeCommitTimeCalculator.js
└── tsconfig.json

28 directories, 100 files
```

