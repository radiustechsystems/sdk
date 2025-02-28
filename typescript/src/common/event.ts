import { BytesLike } from '../providers/eth';

/**
 * Event represents an EVM contract event emitted during transaction execution
 * Contains decoded event data and the raw event payload
 */
export class Event {
  /**
   * Creates a new Event with the given name, data, and raw bytes
   * @param name The name of the event
   * @param data The decoded data of the event as key-value pairs
   * @param raw The raw bytes of the event
   */
  constructor(
    /**
     * The name of the event
     */
    public name: string,
    /**
     * The data of the event as key-value pairs
     */
    public data: Record<string, unknown>,
    /**
     * The raw bytes of the event
     */
    public raw: BytesLike
  ) {}
}
