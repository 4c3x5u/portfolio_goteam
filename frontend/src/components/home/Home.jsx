import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import {
  faUsers,
  faChalkboardTeacher,
  faQuestionCircle,
} from '@fortawesome/free-solid-svg-icons';
import ControlButton from './ControlButton';

import './home.sass';
import logo from '../../assets/homeheader.svg';
import inbox from '../../assets/inbox.svg';
import ready from '../../assets/ready.svg';
import go from '../../assets/go.svg';
import done from '../../assets/done.svg';

const Home = () => (
  <div id="Home">
    <div className="Logo">
      <img alt="logo" src={logo} />
    </div>
    <div className="ControlBar">
      <Container>
        <Row className="Controls">
          <ControlButton
            name="team"
            action={() => console.log('team button clicked')}
            icon={faUsers}
          />

          <ControlButton
            name="boards"
            action={() => console.log('boards button clicked')}
            icon={faChalkboardTeacher}
          />

          <ControlButton
            name="help"
            action={() => console.log('boards button clicked')}
            icon={faQuestionCircle}
          />
        </Row>
      </Container>
    </div>

    <Row className="ColumnsRow">
      <Col className="Col" xs={3}>
        <div className="Column InboxColumn">
          <div className="Header">
            <img src={inbox} alt="inbox column header" />
          </div>
          <div className="Body" />
        </div>
      </Col>

      <Col className="Col" xs={3}>
        <div className="Column ReadyColumn">
          <div className="Header">
            <img src={ready} alt="ready column header" />
          </div>
          <div className="Body" />
        </div>
      </Col>

      <Col className="Col" xs={3}>
        <div className="Column GoColumn">
          <div className="Header">
            <img src={go} alt="go column header" />
          </div>
          <div className="Body" />
        </div>
      </Col>

      <Col className="Col" xs={3}>
        <div className="Column DoneColumn">
          <div className="Header">
            <img src={done} alt="done column header" />
          </div>
          <div className="Body" />
        </div>
      </Col>
    </Row>
  </div>
);

export default Home;
