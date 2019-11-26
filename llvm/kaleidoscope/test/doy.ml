open Unix
(* open Core *)

let t = Unix.localtime (Unix.time ());;

Printf.printf "Today is day %d of the current year.\n" t.tm_yday
