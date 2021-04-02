import data from './data.js';
import jsonGraphqlExpress from 'json-graphql-server';
const app = require('express')()

app.use('/api/graphql', jsonGraphqlExpress(data));


const port = process.env.PORT || 3000;

module.exports = app.listen(port, () => console.log(`Server running on ${port}, http://localhost:${port}`));