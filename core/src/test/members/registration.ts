import { app, mockEventStreamWebSocket } from '../common';
import nock from 'nock';
import request from 'supertest';
import assert from 'assert';
import { IEventMemberRegistered, IDBMember } from '../../lib/interfaces';
import * as utils from '../../lib/utils';
import { timeStamp } from 'console';

describe('Members - registration', async () => {

  const timestampCreation = utils.getTimestamp();
  const timestampUpdate = utils.getTimestamp();

  it('Checks that adding a member sends a request to API Gateway and updates the database', async () => {

    nock('https://apigateway.kaleido.io')
      .post('/registerMember?kld-from=0x0000000000000000000000000000000000000011&kld-sync=true')
      .reply(200);
    const addMemberResponse = await request(app)
      .put('/api/v1/members')
      .send({
        address: '0x0000000000000000000000000000000000000011',
        name: 'Member 1',
        app2appDestination: 'kld://app2app',
        docExchangeDestination: 'kld://docexchange'
      })
      .expect(200);
    assert.deepStrictEqual(addMemberResponse.body, { status: 'submitted' });

    const getMemberResponse = await request(app)
      .get('/api/v1/members')
      .expect(200);
    const member = getMemberResponse.body.find((member: IDBMember) => member.address === '0x0000000000000000000000000000000000000011');
    assert.strictEqual(member.address, '0x0000000000000000000000000000000000000011');
    assert.strictEqual(member.name, 'Member 1');
    assert.strictEqual(member.app2appDestination, 'kld://app2app');
    assert.strictEqual(member.docExchangeDestination, 'kld://docexchange');
    assert.strictEqual(member.owned, true);
    assert.strictEqual(member.confirmed, false);
    assert.strictEqual(typeof member.timestamp, 'number');

    const getMemberByAddressResponse = await request(app)
    .get('/api/v1/members/0x0000000000000000000000000000000000000011')
    .expect(200);
    assert.deepStrictEqual(member, getMemberByAddressResponse.body);
  });

  it('Checks that event stream notification for confirming member registrations is handled', async () => {

    const eventPromise = new Promise((resolve) => {
      mockEventStreamWebSocket.once('send', message => {
        assert.strictEqual(message, '{"type":"ack","topic":"dev"}');
        resolve();
      })
    });

    const data: IEventMemberRegistered = {
      member: '0x0000000000000000000000000000000000000011',
      name: 'Member 1',
      app2appDestination: 'kld://app2app',
      docExchangeDestination: 'kld://docexchange',
      timestamp: timestampCreation
    };
    mockEventStreamWebSocket.emit('message', JSON.stringify([{
      signature: utils.contractEventSignatures.MEMBER_REGISTERED,
      data
    }]));
    await eventPromise;
  });

  it('Get member should return the confirmed member', async () => {
    const getMemberResponse = await request(app)
      .get('/api/v1/members')
      .expect(200);
    const member = getMemberResponse.body.find((member: IDBMember) => member.address === '0x0000000000000000000000000000000000000011');
    assert.strictEqual(member.address, '0x0000000000000000000000000000000000000011');
    assert.strictEqual(member.name, 'Member 1');
    assert.strictEqual(member.app2appDestination, 'kld://app2app');
    assert.strictEqual(member.docExchangeDestination, 'kld://docexchange');
    assert.strictEqual(member.confirmed, true);
    assert.strictEqual(member.owned, true);
    assert.strictEqual(member.timestamp, timestampCreation);

    const getMemberByAddressResponse = await request(app)
    .get('/api/v1/members/0x0000000000000000000000000000000000000011')
    .expect(200);
    assert.deepStrictEqual(member, getMemberByAddressResponse.body);
  });

  it('Checks that updating a member sends a request to API Gateway and updates the database', async () => {
    nock('https://apigateway.kaleido.io')
      .post('/registerMember?kld-from=0x0000000000000000000000000000000000000011&kld-sync=true')
      .reply(200);
    const addMemberResponse = await request(app)
      .put('/api/v1/members')
      .send({
        address: '0x0000000000000000000000000000000000000011',
        name: 'Member 2',
        app2appDestination: 'kld://app2app2',
        docExchangeDestination: 'kld://docexchange2'
      })
      .expect(200);
    assert.deepStrictEqual(addMemberResponse.body, { status: 'submitted' });

    const getMemberResponse = await request(app)
      .get('/api/v1/members')
      .expect(200);
    const member = getMemberResponse.body.find((member: IDBMember) => member.address === '0x0000000000000000000000000000000000000011');
    assert.strictEqual(member.address, '0x0000000000000000000000000000000000000011');
    assert.strictEqual(member.name, 'Member 2');
    assert.strictEqual(member.app2appDestination, 'kld://app2app2');
    assert.strictEqual(member.docExchangeDestination, 'kld://docexchange2');
    assert.strictEqual(member.owned, true);
    assert.strictEqual(member.confirmed, false);
    assert.strictEqual(typeof member.timestamp, 'number');

    const getMemberByAddressResponse = await request(app)
    .get('/api/v1/members/0x0000000000000000000000000000000000000011')
    .expect(200);
    assert.deepStrictEqual(member, getMemberByAddressResponse.body);
  });

  it('Checks that event stream notification for confirming member registrations are handled', async () => {

    const eventPromise = new Promise((resolve) => {
      mockEventStreamWebSocket.once('send', message => {
        assert.strictEqual(message, '{"type":"ack","topic":"dev"}');
        resolve();
      })
    });

    const data: IEventMemberRegistered = {
      member: '0x0000000000000000000000000000000000000011',
      name: 'Member 2',
      app2appDestination: 'kld://app2app2',
      docExchangeDestination: 'kld://docexchange2',
      timestamp: timestampUpdate
    };
    mockEventStreamWebSocket.emit('message', JSON.stringify([{
      signature: utils.contractEventSignatures.MEMBER_REGISTERED,
      data
    }]));
    await eventPromise;
  });

  it('Get member should return the confirmed member', async () => {
    const getMemberResponse = await request(app)
      .get('/api/v1/members')
      .expect(200);
    const member = getMemberResponse.body.find((member: IDBMember) => member.address === '0x0000000000000000000000000000000000000011');
    assert.strictEqual(member.address, '0x0000000000000000000000000000000000000011');
    assert.strictEqual(member.name, 'Member 2');
    assert.strictEqual(member.app2appDestination, 'kld://app2app2');
    assert.strictEqual(member.docExchangeDestination, 'kld://docexchange2');
    assert.strictEqual(member.owned, true);
    assert.strictEqual(member.confirmed, true);
    assert.strictEqual(member.timestamp, timestampUpdate);

    const getMemberByAddressResponse = await request(app)
    .get('/api/v1/members/0x0000000000000000000000000000000000000011')
    .expect(200);
    assert.deepStrictEqual(member, getMemberByAddressResponse.body);
  });

});