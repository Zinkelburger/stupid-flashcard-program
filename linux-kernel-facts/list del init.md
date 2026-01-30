Why use list_del_init instead of just deleting it? If you use a standard delete, the node still contains pointers to the list it was just removed from. By re-initializing it, you ensure that:

    The node is in a "clean" state.

    You can check list_empty(entry) on the deleted node later without causing a crash.

    The node is ready to be added to a different list immediately.

  /**
 * list_del_init - deletes entry from list and reinitialize it.
 * @entry: the element to delete from the list.
 */
static inline void list_del_init(struct list_head *entry)
{
	__list_del_entry(entry);
	INIT_LIST_HEAD(entry);
}