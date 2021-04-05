import React, { useState } from 'react';
import { Container, Row } from 'react-bootstrap';
import {
  faChalkboardTeacher, faQuestionCircle, faUsers,
} from '@fortawesome/free-solid-svg-icons';

import ControlButton from './ControlButton';

import './appheader.sass';
import logo from '../../../assets/homeHeader.svg';

const AppHeader = () => {
  const [teamCtrlIsToggled, setTeamCtrlIsToggled] = useState(false);

  return (
    <div id="AppHeader">
      <div className="Logo">
        <img alt="logo" src={logo} />
      </div>
      <div className="ControlBar">
        <Container>
          <Row className="Controls">
            <ControlButton
              name="team"
              setIsToggled={() => setTeamCtrlIsToggled(!teamCtrlIsToggled)}
              isToggled={teamCtrlIsToggled}
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
    </div>
  );
};

export default AppHeader;
