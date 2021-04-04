import React from 'react';
import { Col, Container, Row } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faGithub, faLinkedin } from '@fortawesome/free-brands-svg-icons';

const Footer = () => (
  <div className="Footer">
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

export default Footer;
