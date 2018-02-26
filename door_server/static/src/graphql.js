import PropTypes from "prop-types";

function makeGraphQLRequest(query, variables, callback) {
    const xhr = new XMLHttpRequest();
    xhr.responseType = 'json';
    xhr.open("POST", "http://localhost:3001/graphql");
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.setRequestHeader("Accept", "application/json");
    xhr.onload = function () {
        callback(xhr.response)
    };
    xhr.send(JSON.stringify({
        query: query,
        variables: variables
    }));
}

makeGraphQLRequest.propTypes = {
    query: PropTypes.string.isRequired,
    variables: PropTypes.object.isRequired,
    callback: PropTypes.func.isRequired
};

export default makeGraphQLRequest;