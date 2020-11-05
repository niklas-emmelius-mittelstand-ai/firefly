import * as database from '../clients/database';
import * as apiGateway from '../clients/api-gateway';
import * as utils from '../lib/utils';
import { IEventMemberRegistered } from '../lib/interfaces';
import RequestError from '../lib/request-error';

export const handleGetMembersRequest = (skip: number, limit: number, owned: boolean) => {
  return database.retrieveMembers(skip, limit, owned);
};

export const handleGetMemberRequest = async (address: string) => {
  const member = await database.retrieveMemberByAddress(address);
  if(member === null) {
    throw new RequestError('Member not found', 404);
  }
  return member;
};

export const handleUpsertMemberRequest = async (address: string, name: string,
  app2appDestination: string, docExchangeDestination: string) => {
  await apiGateway.upsertMember(address, name, app2appDestination, docExchangeDestination);
  await database.upsertMember(address, name, app2appDestination, docExchangeDestination, utils.getTimestamp(), false, true);
};

export const handleMemberRegisteredEvent = async ({ member, name, app2appDestination, docExchangeDestination, timestamp }: IEventMemberRegistered) => {
  const owner = await database.retrieveMemberByAddress(member);
  if(owner !== null && owner.owned) {
    await database.confirmMember(member, Number(timestamp));
  } else {
    await database.upsertMember(member, name, app2appDestination, docExchangeDestination, Number(timestamp), true, false);
  }
};