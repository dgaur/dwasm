# TODO

Laundry list of bugs, incomplete bits, etc:
* Not all section types are decoded/recognized
* Custom 'name' sections can be further decoded and used to annotate functions,
  modules, etc.  See appendix 7.4
* Integrate a better logging module/support: levels, multiple threads, etc
* Add module/section validation and make -v option meaningful
* Add a formal `trap()` path in the VM for runtime exceptions
* Not sure the `end` and `ret` semantics are correct.  Possible that these are
  mixed up incorrectly.  Or possible that the sample code is wrong.  etc.
* Function locals aren't properly cleaned up on function return/exit.



