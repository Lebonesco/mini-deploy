import React from 'react';
import { Input, Divider, Transition, Button } from 'semantic-ui-react';

class Home extends React.Component {
    constructor(props) {
        super(props);
        this.state = { repo: "react-cosmos/create-react-app-example.git", port: "", "loading": false }   
        
        this.deployApp = this.deployApp.bind(this);
    }

    deployApp() {
        const { repo, port } = this.state;
        this.setState({ loading: true })
        fetch(`http://localhost:8080/deploy?repo=https://github.com/${repo}&port=${port}`, {
          method: 'GET'
        })
        .then(res => {
          if (!res.ok) { throw res }
          return res.json()
        })
        .then(data => {
          this.setState({ loading: false })
          let port = data.port
          window.open("http://localhost:" + port, '_blank');
        })
        .catch(err => {
          this.setState({ loading: false })
          console.log(err)
        })
      }

    handleRepoChange = (event, data) => this.setState({ repo: data.value })
    handlePortChange = (event, data) => this.setState({ port: data.value })

    render() {
        return (
            <div>
                <p>
                Paste your <b>Github</b> repo url here
                </p>
                <Input value={this.state.repo} style={{minWidth: '500px'}} inverted label="https://github.com/" placeholder="repo/path.git" size="small" onChange={this.handleRepoChange}/>
                <Divider hidden />
                <Transition.Group animation="slide down" duration={500} style={{marginBottom: "20px"}}>
                { (this.state.repo.length > 0) && 
                    <div>
                    <Input label="port" inverted placeholder="default to 8080" onChange={this.handlePortChange} size="mini"/>
                    <Divider hidden />
                    <Button color="green" size="large" loading={this.state.loading} onClick={this.deployApp}>Deploy</Button>
                    </div>
                }
                </Transition.Group>
            </div>
        )
    }
}

export default Home;