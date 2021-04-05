import React from 'react';
import { Col, Container, Row } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faGithub, faLinkedin } from '@fortawesome/free-brands-svg-icons';

import './appfooter.sass';

const AppFooter = () => (
  <div id="AppFooter">
    <Container>
      <Row className="Row">
        <Col className="Col" xs={6}>
          <a href="https:FIX THIS">
            <FontAwesomeIcon icon={faGithub} />
            PROJECT
          </a>
        </Col>
        <Col className="Col" xs={6}>
          <a href="https:FIX THIS">
            <FontAwesomeIcon icon={faLinkedin} />
            AUTHOR
          </a>
        </Col>
      </Row>
    </Container>
  </div>
);

export default AppFooter;
