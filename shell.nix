(import
  (
    let
      lock = builtins.fromJSON (builtins.readFile ./flake.lock);
    in
      fetchTarball (with lock.nodes.flake-compat.locked; {
        url = "https://github.com/${owner}/${repo}/archive/${rev}.tar.gz";
        sha256 = narHash;
      })
  )
  {
    src = ./.;
  })
.shellNix
