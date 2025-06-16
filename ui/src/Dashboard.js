import React, { useEffect, useRef } from 'react';
import { styled, ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import GitHubIcon from '@mui/icons-material/GitHub';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';
import Link from '@mui/material/Link';

const theme = createTheme();

function Source() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      <Link color="inherit" href="https://github.com/n3wscott/colorgrid">
        <GitHubIcon/>
      </Link>
    </Typography>
  );
}


const Root = styled('div')(({ theme }) => ({
  display: 'flex',
}));

const StyledToolbar = styled(Toolbar)(({ theme }) => ({
  paddingRight: 24,
}));

const StyledAppBar = styled(AppBar)(({ theme }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  backgroundImage: 'linear-gradient(to right, red, orange, yellow, green, cyan, blue, violet)',
}));

const StyledTitle = styled(Typography)(({ theme }) => ({
  flexGrow: 1,
}));

const AppBarSpacer = styled('div')(({ theme }) => ({
  ...theme.mixins.toolbar,
}));

const Content = styled('main')(({ theme }) => ({
  flexGrow: 1,
  height: '100vh',
  overflow: 'auto',
}));

const StyledContainer = styled(Container)(({ theme }) => ({
  paddingTop: theme.spacing(4),
  paddingBottom: theme.spacing(4),
}));

const ClientFrame = styled('iframe')(({ theme }) => ({
  height: 100,
  width: 100,
}));


export const useInterval = (callback, delay) => {
  const savedCallback = useRef();
  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);
  useEffect(() => {
    function tick() {
      savedCallback.current();
    }
    if (delay !== null) {
      const id = setInterval(tick, delay);
      return () => clearInterval(id);
    }
  }, [delay]);
}

function ClientBox(props) {
  const [count, setCount] = React.useState(0);

  let interval = 3800 + Math.floor(Math.random() * 500);
  useInterval(() => {
      setCount(count+1);
    }, interval);

  return (
      <ClientFrame title={props.key} key={count} src={props.url} />
  );
}

function useQuery() {
  return new URLSearchParams(window.location.search);
}

export default function Dashboard(props) {
  let query = useQuery();
  const url = query.get("url");
  let count = query.get("count");
  console.log(count)
  if (count === null) {
    count = 10;
  } else {
    count = parseInt(count)
  }
  const items = new Array(count).fill(null);

  return (
    <ThemeProvider theme={theme}>
      <Root>
        <CssBaseline />
        <StyledAppBar position="absolute" color="transparent">
          <StyledToolbar>
            <StyledTitle component="h1" variant="h6" color="inherit" noWrap>
              ColorGrid, Rainbow Deployments Viewer
            </StyledTitle>
          </StyledToolbar>
        </StyledAppBar>
        <Content>
          <AppBarSpacer />

          <StyledContainer maxWidth="lg">
            <Paper>
              {items.map((_, idx) => <ClientBox key={idx} url={url}/>)}
            </Paper>
            <Box pt={4}>
              <Source />
            </Box>
          </StyledContainer>
        </Content>
      </Root>
    </ThemeProvider>
  );
}
