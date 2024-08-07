# ProtonDB Community API

API for developers seeking to leverage data from ProtonDB.
##### Note:This is an deployment with sample data
https://protondb-community-api-04f42bc1742f.herokuapp.com/

## Introduction

The ProtonDB Community API is an open-source project designed to provide developers with programmatic access to compatibility data of Windows games running on Linux through Proton. Proton, developed by Valve Corporation, is a compatibility layer that enables gaming on Linux systems by using technologies like Wine, DXVK, and VKD3D. ProtonDB, on the other hand, is a community-driven database that collects user reports on game compatibility with Proton.

This API serves as a valuable resource for developers looking to build applications, tools, and services that utilize ProtonDB's gaming compatibility information. It acts as a bridge between developers and the ProtonDB database, allowing seamless access to real-world feedback on how well specific games perform on Linux.

## Features

- **Automatic Data Updates:** The API automatically checks and adds the latest data dump every 31 days, ensuring that you have access to the most up-to-date compatibility information.

- **Game-Specific Filtering:** Developers can filter compatibility reports based on specific games, allowing them to retrieve data for individual titles.

- **Versioned Data Structure:** For reports inserted before December 2019, the API provides access to versioned data structures, ensuring compatibility with historical reports and analysis.

- **Game Summary Access:** Developers can access a game's summary, including tiers fetched directly from ProtonDB, providing essential information for game performance assessment.

- **Stats Endpoint:** The API provides a `/api/stats` endpoint that allows developers to retrieve statistics about the API usage. It includes information on the number of requests, average response time, and the time remaining for the next automatic data update.

- **Last Processed File:** The API offers the ability to view the last processed data dump file, giving insights into the latest data available.

## Installation

1. Clone the repository to your local machine:

```bash
git clone https://github.com/trsnaqe/protondb-community-api.git
```

2. Install MongoDB on your system and set up a local MongoDB database.

3. Create a `.env` file in the root directory of the project and set the MongoDB connection URI:

```bash
DB_URI=mongodb+srv://<username>:<password>@cluster0.xyz.mongodb.net/?retryWrites=true&w=majority
```

4. Open a terminal or command prompt and navigate to the project's directory.

5. Run the project using the `go run` command:

The API will now be up and running, and you can start making requests to the available endpoints. Ensure that MongoDB is running and accessible via the connection URI specified in the `.env` file.

## API Documentation

- `/api/games (GET)`: Get all games. [Disabled: The dataset is large and costs a lot to leave this endpoint open.]

- `/api/games/{gameId} (GET)`: Get a game by gameId.
- `/api/games/{gameId}/summary (GET)`: Get tiers by gameId, fetched from ProtonDB directly.

- `/api/reports (GET)`: Retrieve reports; add `?versioned=true` for versioned data. [Disabled: The dataset is large and costs a lot to leave this endpoint open.]

- `/api/reports/{gameId} (GET)`: Get reports by gameId; add `?versioned=true` for versioned data.

- `/api/stats (GET)`: Get stats of the API. This endpoint provides information about API usage, response times, and the time remaining for the next automatic data update.

- `/api/v2/games`: Get games endpoint. Supports query in v2. If no query is present gets all the games. Query options: [gameid|game_id|appid|app_id] to get game with id. [title] to search game by title, use [precision] to increase or decrease matching accuracy, value should be higher than 0. If both title and appid are present, appid supersedes.

- `/api/v2/reports`: Get reports endpoint. Supports query in v2. If no query is present gets all the reports. Query options: [gameid|game_id|appid|app_id] to get game with id. [title] to search game by title, use [precision] to increase or decrease matching accuracy, value should be higher than 0. If both title and appid are present, appid supersedes. [versioned] to get the reports with metadata. [version] 1 or 2 to filter by report versions.

## Contributing

We welcome contributions to the project! Whether you want to report issues, submit feature requests, or make pull requests, your input is valuable in improving the Linux gaming experience. Please refer to our [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to contribute.

## License

This project is open-source and available under the [MIT License](LICENSE).

## Acknowledgments

We would like to acknowledge the following for their contributions and inspiration to the ProtonDB Community API:

- [bdefore](https://github.com/bdefore)
- [ProtonDB](https://protondb.com/)
- [Proton from Valve](https://github.com/ValveSoftware/Proton)

## Contact

For questions or support, you can reach out via the following channels:

- [LinkedIn](https://www.linkedin.com/in/sacit)
- [Twitter](https://twitter.com/Trsnaqe)
- [Email](trsnaqe@gmail.com)
- [Buy Me a Coffee](https://www.buymeacoffee.com/trsnaqe)

We appreciate your contribution and look forward to making the ProtonDB Community API even better together!

### Important Update

As of 04/09/2023 I had to terminate the deployment due to budget reasons, as API wasn't getting enough traffic to justify its deployment cost. This service is ready to deploy and it can be used locally; 

### How to Host

1. Clone the repository to your local machine:

```bash
git clone https://github.com/trsnaqe/protondb-community-api.git
```

2. Install MongoDB on your system and set up a local MongoDB database.

3. Create a `.env` file in the root directory of the project and set the MongoDB connection URI:

```bash
DB_URI=mongodb+srv://<username>:<password>@cluster0.xyz.mongodb.net/?retryWrites=true&w=majority
```
4. Download database backup. (It contains everything up to August 3th of 2023)

https://drive.google.com/file/d/1tE8T7TRjxYlZpw2qKaKWv9NneR6m7dP8/view?usp=sharing

5. Open a terminal or command prompt and navigate to the Downloaded file's directory.

6. Run the below command

```bash
sudo mongorestore -db --uri "mongodb+srv://<url>&w=majority" --dir .
```
